package database

import (
	"bikesRentalAPI/internal/users/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSeeder(t *testing.T) {
	// GIVEN: a database base url
	t.Setenv("DB_URL", "file::memory:?cache=shared")

	// GIVEN: a database service
	dbService := New()
	dbService.Start()
	defer dbService.Close()

	// WHEN: the seeder is created
	seeder := NewSeeder(dbService)

	// THEN: Assert that the seeder is not nil and Seeder interface is implemented
	assert.NotNil(t, seeder)
	assert.Implementsf(t, (*Seeder)(nil), seeder, "seeder does not implement Seeder")
}

func TestSeed(t *testing.T) {
	testValues := []struct {
		name          string
		creds         string
		expectedError error
	}{
		{
			name:          "success - admin user seeded: 'user@email.com:password'",
			creds:         "'user@email.com:password",
			expectedError: nil,
		},
		{
			name:          "failure - admin credentials wrongly encoded",
			creds:         "admin.admin",
			expectedError: fmt.Errorf("failed Seed admin: failed to decode admin credentials."),
		},
	}
	for test := range testValues {
		t.Run(testValues[test].name, func(t *testing.T) {
			// GIVEN: a database base url and admin credentials
			t.Setenv("DB_URL", "file::memory:?cache=shared")
			t.Setenv("USER_CREDENTIALS", testValues[test].creds)

			// GIVEN: a database service
			dbService := New()
			dbService.Start()
			defer dbService.Close()

			// GIVEN: a seeder
			seeder := NewSeeder(dbService)

			// WHEN: calls to Seed
			err := seeder.Seed(models.User{})

			if testValues[test].expectedError != nil {

				// THEN: Assert that the seeder is not nil
				assert.NotNil(t, seeder)
				assert.ErrorAs(t, err, &testValues[test].expectedError)

			} else {
				// THEN: Assert that the seeder is not nil and no error is returned
				assert.NotNil(t, seeder)
				assert.NoError(t, err)
			}
		})
	}
}
