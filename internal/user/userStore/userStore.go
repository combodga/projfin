package userStore

import (
	"context"
	"fmt"

	"github.com/combodga/projfin/internal/store"
	"github.com/lib/pq"
)

func DoRegister(s *store.Store, ctx context.Context, username, password string) error {
	sql := "INSERT INTO users VALUES ($1, $2, 0, 0)"
	_, err := s.DB.ExecContext(ctx, sql, username, password)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return s.ErrorDupe
			}
		}
		return fmt.Errorf("store query error: %w", err)
	}
	return nil
}

func DoLogin(s *store.Store, ctx context.Context, username, password string) error {
	var count int
	sql := "SELECT COUNT(*) FROM users WHERE username = $1 AND password = $2"
	rows, err := s.DB.QueryxContext(ctx, sql, username, password)
	if err != nil {
		return fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return fmt.Errorf("store query rows error: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("auth error: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("store get rows error: %w", err)
	}

	return nil
}
