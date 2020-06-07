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

func doSetupDB(dbfile string, username string, password string) (*sql.DB, error) {
	db, err := CreateSqliteDBFile(dbfile, true)
	if err != nil {
		return nil, err
	}
	err = sqly.ExecuteQueryFromFile(db, "user.sql")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`INSERT INTO users(username, pwhash) VALUES(? , ?) `, username, password)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestSelectRowLegacy(t *testing.T) {
	username := "user1"
	password := "password"

	db, err := doSetupDB("selectrow.db", username, password)
	require.Nil(t, err)

	var storedPassword string
	err = sqly.SelectRowLegacy(db,
		func(row *sql.Row) error {
			return row.Scan(&storedPassword)
		},
		`SELECT pwhash from users where username = ?`,
		username,
	)
	require.Nil(t, err)

	assert.Equal(t, password, storedPassword)
}

func TestSelectRow(t *testing.T) {
	username := "user1"
	password := "password"

	db, err := doSetupDB("selectrow2.db", username, password)
	require.Nil(t, err)

	var storedPassword string
	err = sqly.SelectRow(db,
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

	db, err := doSetupDB("select.db", username, password)
	require.Nil(t, err)

	username2 := "user2"
	_, err = db.Exec(`INSERT INTO users(username, pwhash) VALUES(? , ?) `, username2, password)
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

func BenchmarkSelectRowLegacy(b *testing.B) {
	username := "user1"
	password := "password"

	db, err := doSetupDB("selectrow.db", username, password)
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		var storedPassword string
		err = sqly.SelectRowLegacy(db,
			func(row *sql.Row) error {
				return row.Scan(&storedPassword)
			},
			`SELECT pwhash from users where username = ?`,
			username,
		)
	}
}

func BenchmarkSelectRow(b *testing.B) {
	username := "user1"
	password := "password"

	db, err := doSetupDB("selectrow.db", username, password)
	if err != nil {
		return
	}
	for n := 0; n < b.N; n++ {
		var storedPassword string
		err = sqly.SelectRow(db,
			func(row *sql.Row) error {
				return row.Scan(&storedPassword)
			},
			`SELECT pwhash from users where username = ?`,
			username,
		)
	}
}
