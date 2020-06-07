package sqly_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/lapiscrm/sqly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Import sqlite3 driver here
	_ "github.com/mattn/go-sqlite3"
)

func CreateSqliteDBFile(dbfile string, deleteIfExists bool) (*sql.DB, error) {
	if deleteIfExists {
		if _, err := os.Stat(dbfile); err == nil {
			err := os.Remove(dbfile)
			if err != nil {
				return nil, err
			}
		} else if os.IsNotExist(err) {
			// No such file
		} else {
			return nil, err
		}
	}
	return sql.Open("sqlite3", dbfile)
}

func doSetupDB(t *testing.T, dbfile string, username string, password string) *sql.DB {
	db, err := CreateSqliteDBFile(dbfile, true)
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

func TestSelect(t *testing.T) {
	username := "user1"
	password := "password"

	db := doSetupDB(t, "selectrow2.db", username, password)
	username2 := "user2"
	_, err := db.Exec(`INSERT INTO users(username, pwhash) VALUES(? , ?) `, username2, password)
	require.Nil(t, err)

	users := []string{}
	err = sqly.Select(db,
		func(rows *sql.Rows) error {
			var storedUsername string
			err := rows.Scan(&storedUsername)
			if err != nil {
				return err
			}
			users = append(users, storedUsername)
			return nil
		},
		`SELECT username from users`,
	)
	require.Nil(t, err)
	assert.Equal(t, len(users), 2)
}
