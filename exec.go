package sqly

import (
	"context"
	"database/sql"
	"io/ioutil"
)

func ExecuteQueryFromFileContext(ctx context.Context, DB *sql.DB, queryFile string) error {
	dat, err := ioutil.ReadFile(queryFile)
	if err != nil {
		return err
	}
	_, err = DB.ExecContext(ctx, string(dat))
	return err
}

func ExecuteQueryFromFile(DB *sql.DB, queryFile string) error {
	return ExecuteQueryFromFileContext(context.Background(), DB, queryFile)
}

func ExecuteQueryFromFileTxContext(ctx context.Context, tx *sql.Tx, queryFile string) error {
	dat, err := ioutil.ReadFile(queryFile)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, string(dat))
	return err
}

func ExecuteQueryFromFileTx(tx *sql.Tx, queryFile string) error {
	return ExecuteQueryFromFileTxContext(context.Background(), tx, queryFile)
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
