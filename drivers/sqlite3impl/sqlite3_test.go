package sqlite3impl_test

import (
	"testing"

	"github.com/lapiscrm/sqly/drivers/sqlite3impl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSqliteDBFile(t *testing.T) {
	db, err := sqlite3impl.CreateSqliteDBFile("test.db", true, nil)
	require.Nil(t, err)
	assert.NotNil(t, db)
}
