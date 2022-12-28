package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	DB        *sqlx.DB
	ErrorDupe error
}

func New(database string) (*Store, error) {
	db, err := sqlx.Connect("postgres", database)
	if err != nil {
		return nil, fmt.Errorf("store connect error: %w", err)
	}

	s := &Store{
		DB:        db,
		ErrorDupe: fmt.Errorf("duplicate key error"),
	}

	sql := "BEGIN;CREATE TABLE IF NOT EXISTS users (username text primary key,password text,balance double precision,withdrawn double precision);"
	sql += "CREATE TABLE IF NOT EXISTS orders (order_number text primary key,username text REFERENCES users,status text,accrual double precision,uploaded_at timestamp with time zone);"
	sql += "CREATE TABLE IF NOT EXISTS withdrawals (order_number text,username text REFERENCES users,sum double precision,processed_at timestamp with time zone);COMMIT;"
	db.MustExec(sql)

	return s, nil
}

func (s *Store) Close() {
	s.DB.Close()
}
