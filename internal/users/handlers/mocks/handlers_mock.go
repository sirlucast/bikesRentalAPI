// Code generated by MockGen. DO NOT EDIT.
// Source: internal/users/handlers/handlers.go
//
// Generated by this command:
//
//	mockgen -source=internal/users/handlers/handlers.go -destination=internal/users/handlers/mocks/handlers_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	jwtauth "github.com/go-chi/jwtauth/v5"
	gomock "go.uber.org/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// GetUserDetails mocks base method.
func (m *MockHandler) GetUserDetails(w http.ResponseWriter, req *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetUserDetails", w, req)
}

// GetUserDetails indicates an expected call of GetUserDetails.
func (mr *MockHandlerMockRecorder) GetUserDetails(w, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDetails", reflect.TypeOf((*MockHandler)(nil).GetUserDetails), w, req)
}

// GetUserProfile mocks base method.
func (m *MockHandler) GetUserProfile(w http.ResponseWriter, req *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetUserProfile", w, req)
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockHandlerMockRecorder) GetUserProfile(w, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockHandler)(nil).GetUserProfile), w, req)
}

// ListAllUsers mocks base method.
func (m *MockHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListAllUsers", w, r)
}

// ListAllUsers indicates an expected call of ListAllUsers.
func (mr *MockHandlerMockRecorder) ListAllUsers(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllUsers", reflect.TypeOf((*MockHandler)(nil).ListAllUsers), w, r)
}

// LoginUser mocks base method.
func (m *MockHandler) LoginUser(tokenAuth *jwtauth.JWTAuth, w http.ResponseWriter, req *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LoginUser", tokenAuth, w, req)
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockHandlerMockRecorder) LoginUser(tokenAuth, w, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockHandler)(nil).LoginUser), tokenAuth, w, req)
}

// RegisterUser mocks base method.
func (m *MockHandler) RegisterUser(w http.ResponseWriter, req *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterUser", w, req)
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockHandlerMockRecorder) RegisterUser(w, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockHandler)(nil).RegisterUser), w, req)
}

// UpdateUserDetails mocks base method.
func (m *MockHandler) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateUserDetails", w, r)
}

// UpdateUserDetails indicates an expected call of UpdateUserDetails.
func (mr *MockHandlerMockRecorder) UpdateUserDetails(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserDetails", reflect.TypeOf((*MockHandler)(nil).UpdateUserDetails), w, r)
}

// UpdateUserProfile mocks base method.
func (m *MockHandler) UpdateUserProfile(w http.ResponseWriter, req *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateUserProfile", w, req)
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockHandlerMockRecorder) UpdateUserProfile(w, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockHandler)(nil).UpdateUserProfile), w, req)
}
