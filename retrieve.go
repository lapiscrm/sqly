package sqly

import (
	"context"
	"database/sql"
)

type onRowFunc func(*sql.Row) error

type onRowsFunc func(*sql.Rows) error

// SelectRowLegacy is a bit slower compared to SelectRow . See the benchmarkets in retrieve_test.go to understand better
func SelectRowLegacy(DB *sql.DB, onRow onRowFunc, query string, args ...interface{}) error {
	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(args...)
	err = onRow(row)
	return err
}

func SelectRowContext(ctx context.Context, DB *sql.DB, onRow onRowFunc, query string, args ...interface{}) error {
	return onRow(DB.QueryRowContext(ctx, query, args...))
}

func SelectRowTxContext(ctx context.Context, tx *sql.Tx, onRow onRowFunc, query string, args ...interface{}) error {
	return onRow(tx.QueryRowContext(ctx, query, args...))
}

func SelectRow(DB *sql.DB, onRow onRowFunc, query string, args ...interface{}) error {
	return SelectRowContext(context.Background(), DB, onRow, query, args...)
}

func SelectRowTx(tx *sql.Tx, onRow onRowFunc, query string, args ...interface{}) error {
	return SelectRowTxContext(context.Background(), tx, onRow, query, args...)
}

func doSelectStmt(ctx context.Context, stmt *sql.Stmt, onRows onRowsFunc, args ...interface{}) error {
	rows, err := stmt.QueryContext(ctx, args...)
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

func SelectTxContext(ctx context.Context, tx *sql.Tx, onRows onRowsFunc, query string, args ...interface{}) error {
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return doSelectStmt(ctx, stmt, onRows, args...)
}

func SelectContext(ctx context.Context, DB *sql.DB, onRows onRowsFunc, query string, args ...interface{}) error {
	stmt, err := DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return doSelectStmt(ctx, stmt, onRows, args...)
}

func Select(DB *sql.DB, onRows onRowsFunc, query string, args ...interface{}) error {
	return SelectContext(context.Background(), DB, onRows, query, args...)
}

func SelectTx(tx *sql.Tx, onRows onRowsFunc, query string, args ...interface{}) error {
	return SelectTxContext(context.Background(), tx, onRows, query, args...)
}
