package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

/*func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresDB(cfg string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}*/

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

/*

	ed := fmt.Errorf("duplicate key error")
	s := &Store{
		DB:        db,
		ErrorDupe: ed,
		Order:     NewOrderPG(db, ed),
		User:      NewUserPG(db, ed),
		Withdraw:  NewWithdrawPG(db),
	}

*/
