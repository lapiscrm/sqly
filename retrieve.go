package sqly

import (
	"database/sql"
)

type onRowFunc func(*sql.Row) error

type onRowsFunc func(*sql.Rows) error

func SelectRow(DB *sql.DB, onRow onRowFunc, query string, args ...interface{}) error {
	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(args...)
	err = onRow(row)
	return err
}

func SelectRow2(DB *sql.DB, onRow onRowFunc, query string, args ...interface{}) error {
	row := DB.QueryRow(query, args...)
	return onRow(row)
}

func Select(DB *sql.DB, onRows onRowsFunc, query string, args ...interface{}) error {
	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := onRows(rows)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}