package sqlite3impl

import (
	"database/sql"
	"os" 

	"github.com/lapiscrm/sqly"

	// Import sqlite3 driver here
	_ "github.com/mattn/go-sqlite3"
)

func CreateSqliteDBFile(dbfile string, deleteIfExists bool, queryFiles []string) (*sql.DB, error) {
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
	DB, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	if queryFiles == nil {
		return DB, nil

	}
	for _, queryFile := range queryFiles {
		err := sqly.ExecuteQueryFromFile(DB, queryFile)
		if err != nil {
			return nil, err
		}
	}
	return DB, nil
}
