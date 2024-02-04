package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	// GIVEN: a database base url
	t.Setenv("DB_URL", "file::memory:?cache=shared")

	// GIVEN: a database service
	dbService, err := New()
	defer dbService.(*database).db.Close()
	assert.NoError(t, err)

	// WHEN: the health check is called
	err = dbService.Health()

	// THEN: Assert that there is no error (indicating a successful health check)
	assert.NoError(t, err)
}

func TestNew(t *testing.T) {
	// GIVEN: a database base url
	os.Setenv("DB_URL", "file::memory:?cache=shared")

	// WHEN: the New function is called
	dbService, err := New()
	defer dbService.(*database).db.Close()
	assert.NoError(t, err)

	// THEN: Assert that the service is not nil
	assert.NotNil(t, dbService)

	// THEN: Assert that the service has a non-nil database connection
	assert.NotNil(t, dbService.(*database).db)
}
