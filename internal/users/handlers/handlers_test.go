package handlers

import (
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/users/models"
	"bikesRentalAPI/internal/users/repository/mocks"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	testEmail      = "test@test.com"
	testPsw        = "test"
	testInvalidPsw = "invalid"
	testSecretKey  = "test_key"
)

var (
	hashedPasw, _   = helpers.GetHashPassword(testPsw)
	mockedValidUser = models.User{
		ID:             1,
		Email:          testEmail,
		HashedPassword: hashedPasw,
	}
	testTokenAuth = &jwtauth.JWTAuth{}
)

func TestLoginUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsersRepo := mocks.NewMockUserRepository(mockCtrl)

	testCases := []struct {
		testJWTAlg          string
		name                string
		email               string
		password            string
		callMock            bool
		mockedUser          models.User
		expectedRepoError   error
		expectedHttpCode    int
		expectedResponseMsg string
	}{
		{
			name:                "Success - LoginUser receives a tokenAuth and a request and returns a response",
			testJWTAlg:          "HS256",
			email:               testEmail,
			password:            testPsw,
			callMock:            true,
			mockedUser:          mockedValidUser,
			expectedRepoError:   nil,
			expectedHttpCode:    http.StatusOK,
			expectedResponseMsg: "Token",
		},
		{
			name:                "Failure - LoginUser receives a tokenAuth and a empty email/password request. Returns error 400",
			testJWTAlg:          "HS256",
			email:               "",
			password:            "",
			callMock:            false,
			mockedUser:          models.User{},
			expectedRepoError:   nil,
			expectedHttpCode:    http.StatusBadRequest,
			expectedResponseMsg: "Validation errors",
		},
		{
			name:                "Failure - LoginUser receives a tokenAuth and a invalid password request. Returns error 401",
			testJWTAlg:          "HS256",
			email:               testEmail,
			password:            testInvalidPsw,
			callMock:            true,
			mockedUser:          mockedValidUser,
			expectedRepoError:   nil,
			expectedHttpCode:    http.StatusUnauthorized,
			expectedResponseMsg: "Invalid username or password",
		},
		{
			name:                "Failure - LoginUser receives a tokenAuth and a non existing email request. Returns error 401",
			testJWTAlg:          "HS256",
			email:               "invalid_email@email.com",
			password:            testPsw,
			callMock:            true,
			mockedUser:          models.User{},
			expectedRepoError:   fmt.Errorf("error"),
			expectedHttpCode:    http.StatusUnauthorized,
			expectedResponseMsg: "Invalid username or password",
		},
		{
			name:                "Failure - LoginUser receives an invalid tokenAuth. Returns error 500",
			testJWTAlg:          "",
			email:               testEmail,
			password:            testPsw,
			callMock:            true,
			mockedUser:          mockedValidUser,
			expectedRepoError:   nil,
			expectedHttpCode:    http.StatusInternalServerError,
			expectedResponseMsg: "Error encoding token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.callMock {
				// GIVEN: a mocked user from repository
				mockUsersRepo.EXPECT().GetUserByEmail(gomock.Any()).Return(tc.mockedUser, tc.expectedRepoError).Times(1)
			}
			// GIVEN: a tokenAuth
			testTokenAuth = jwtauth.New(tc.testJWTAlg, []byte(testSecretKey), nil)
			// GIVEN: a request to login a user
			data := url.Values{}
			data.Set("email", tc.email)
			data.Set("password", tc.password)
			req, err := http.NewRequest("POST", "/users/login", strings.NewReader(data.Encode()))
			assert.Nil(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			// GIVEN a recorder to record the response
			rr := httptest.NewRecorder()
			// GIVEN a user handler
			userHandler := New(mockUsersRepo)
			handler := http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					userHandler.LoginUser(testTokenAuth, w, r)
				},
			)
			// WHEN: the request is made
			handler.ServeHTTP(rr, req)
			// THEN: the status code should be the expected
			assert.Equal(t, tc.expectedHttpCode, rr.Code)
			// THEN: the response should contain the expected error message
			assert.Contains(t, rr.Body.String(), tc.expectedResponseMsg)
		})
	}
}

func TestRegisterUser(t *testing.T) {
	// GIVEN a mocked user repository
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsersRepo := mocks.NewMockUserRepository(mockCtrl)

	testCases := []struct {
		name                string
		createuserReq       models.CreateUserRequest
		mockCreateUser      bool
		mockedErrror        error
		expectedRepoError   error
		expectedHttpCode    int
		expectedResponseMsg string
	}{
		{
			name: "Success - RegisterUser receives a valid request and returns a response",
			createuserReq: models.CreateUserRequest{
				Email:     testEmail,
				Password:  testPsw,
				FirstName: "test",
				LastName:  "test",
			},
			expectedRepoError:   nil,
			mockCreateUser:      true,
			expectedHttpCode:    http.StatusCreated,
			expectedResponseMsg: "1",
		},
		{
			name: "Failure - RegisterUser receives a request with missing values. Handler returns error 400",
			createuserReq: models.CreateUserRequest{
				Email: testEmail,
			},
			expectedRepoError:   nil,
			mockCreateUser:      false,
			expectedHttpCode:    http.StatusBadRequest,
			expectedResponseMsg: "Validation errors",
		},
		{
			name: "Failure - RegisterUser receives a request but Databse returns error. Handler returns error 500",
			createuserReq: models.CreateUserRequest{
				Email:     testEmail,
				Password:  testPsw,
				FirstName: "test",
				LastName:  "test",
			},
			expectedRepoError:   nil,
			mockCreateUser:      true,
			mockedErrror:        assert.AnError,
			expectedHttpCode:    http.StatusInternalServerError,
			expectedResponseMsg: "Error creating user",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the CreateUser method
			if tc.mockCreateUser {
				if tc.mockedErrror != nil {
					mockUsersRepo.EXPECT().CreateUser(gomock.Any()).Return(int64(1), tc.mockedErrror).Times(1)
				} else {
					mockUsersRepo.EXPECT().CreateUser(gomock.Any()).Return(int64(1), nil).Times(1)
				}
			}
			// GIVEN a request to register a user
			body, err := json.Marshal(tc.createuserReq)
			assert.Nil(t, err)
			req, err := http.NewRequest("POST", "/users/register", strings.NewReader(string(body)))
			assert.Nil(t, err)
			// GIVEN a recorder to record the response
			rr := httptest.NewRecorder()

			// GIVEN a user handler
			userHandler := New(mockUsersRepo)
			handler := http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					userHandler.RegisterUser(w, r)
				},
			)
			// WHEN the request is made
			handler.ServeHTTP(rr, req)
			// THEN the user should be registered
			assert.Equal(t, tc.expectedHttpCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedResponseMsg)
		})
	}
}

func TestGetUserProfile(t *testing.T) {
	// GIVEN: a request to get a user's profile
	// WHEN: the request is made
	// THEN: the user's profile should be returned
}

func TestUpdateUserProfile(t *testing.T) {
	// GIVEN: a request to update a user's profile
	// WHEN: the request is made
	// THEN: the user's profile should be updated
}

func TestListUsers(t *testing.T) {
	// GIVEN: a request to list all users
	// WHEN: the request is made
	// THEN: all users should be listed
}

func TestGetUserDetails(t *testing.T) {
	// GIVEN: a request to get user details
	// WHEN: the request is made
	// THEN: the user details should be returned
}

func TestUpdateUserDetails(t *testing.T) {
	// GIVEN: a request to update user details
	// WHEN: the request is made
	// THEN: the user details should be updated
}
