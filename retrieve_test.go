package sqly_test

import (
	"database/sql"
	"testing"

	"github.com/lapiscrm/sqly"
	"github.com/lapiscrm/sqly/drivers/sqlite3impl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func doSetupDB(t *testing.T, dbfile string, username string, password string) *sql.DB {
	db, err := sqlite3impl.CreateSqliteDBFile(dbfile, true, nil)
	require.Nil(t, err)
	require.NotNil(t, db)

	err = sqly.ExecuteQueryFromFile(db, "user.sql")
	require.Nil(t, err)

	_, err = db.Exec(`INSERT INTO users(username, pwhash) VALUES(? , ?) `, username, password)
	require.Nil(t, err)

	return db
}

func TestSelectRow(t *testing.T) {
	username := "user1"
	password := "password"

	db := doSetupDB(t, "selectrow.db", username, password)

	var storedPassword string
	err := sqly.SelectRow(db,
		func(row *sql.Row) error {
			return row.Scan(&storedPassword)
		},
		`SELECT pwhash from users where username = ?`,
		username,
	)
	require.Nil(t, err)

	assert.Equal(t, password, storedPassword)
}

func TestSelectRow2(t *testing.T) {
	username := "user1"
	password := "password"

	db := doSetupDB(t, "selectrow2.db", username, password)

	var storedPassword string
	err := sqly.SelectRow2(db,
		func(row *sql.Row) error {
			return row.Scan(&storedPassword)
		},
		`SELECT pwhash from users where username = ?`,
		username,
	)
	require.Nil(t, err)

	assert.Equal(t, password, storedPassword)
}
