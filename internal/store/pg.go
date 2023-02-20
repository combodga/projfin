package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const PGDuplicateCode = "23505"

func NewPGDB(cfg string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg)
	if err != nil {
		return nil, fmt.Errorf("store DB connect error: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("store DB no connection error: %w", err)
	}

	sql := `
	  BEGIN;
	  CREATE TABLE IF NOT EXISTS users (
		username text primary key,
		password text,
		balance double precision,
		withdrawn double precision
	  );
	  CREATE TABLE IF NOT EXISTS orders (
		order_number text primary key,
		username text REFERENCES users,
		status text,
		accrual double precision,
		uploaded_at timestamp with time zone
	  );
	  CREATE TABLE IF NOT EXISTS withdrawals (
		order_number text,
		username text REFERENCES users,
		sum double precision,
		processed_at timestamp with time zone
	  );
	  COMMIT;`
	db.MustExec(sql)

	return db, nil
}

func PGClose(db *sqlx.DB) {
	db.Close()
}
