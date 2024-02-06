package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// WHEN: the New function is called
	dbService := New()

	// THEN: Assert that the service is not nil
	assert.NotNil(t, dbService)
}

func TestHealth(t *testing.T) {
	// GIVEN: a database base url
	t.Setenv("DB_URL", "file::memory:?cache=shared")

	// GIVEN: a database service
	dbService := New()
	dbService.Start()
	defer dbService.Close()

	// WHEN: the health check is called
	err := dbService.Health()

	// THEN: Assert that there is no error (indicating a successful health check)
	assert.NoError(t, err)
}

func TestStart(t *testing.T) {
	// GIVEN: a database base url
	os.Setenv("DB_URL", "file::memory:?cache=shared")

	// WHEN: the New function is called and Start is called
	dbService := New()
	err := dbService.Start()
	assert.NoError(t, err)
	defer dbService.Close()

	// THEN: Assert that the service is not nil
	assert.NotNil(t, dbService)

	// THEN: Assert that the service has a non-nil database connection
	assert.NotNil(t, dbService.(*database).db)
}

func TestMigrate(t *testing.T) {
	// GIVEN: a database base url
	os.Setenv("DB_URL", "file::memory:?cache=shared")

	// GIVEN: a database service and started
	dbService := New()
	err := dbService.Start()
	assert.NoError(t, err)
	defer dbService.Close()

	// WHEN: the Migrate function is called
	err = dbService.Migrate()

	// THEN: Assert that there is no error (indicating a successful migration)
	assert.NoError(t, err)
}
