package sqly

import (
	"database/sql"
	"io/ioutil"
)

func ExecuteQueryFromFile(DB *sql.DB, queryFile string) error {
	dat, err := ioutil.ReadFile(queryFile)
	if err != nil {
		return err
	}
	_, err = DB.Exec(string(dat))
	return err
}

func ExecuteQueryFromFiles(DB *sql.DB, queryFiles []string) error {
	for _, queryFile := range queryFiles {
		err := ExecuteQueryFromFile(DB, queryFile)
		if err != nil {
			return err
		}
	}
	return nil
}
