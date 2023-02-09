package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserPG_DoRegister(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := NewUserPG(db)

	tests := []struct {
		name     string
		mock     func()
		username string
		password string
		wantErr  bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test", "test").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			username: "test",
			password: "test",
		},
		{
			name: "Empty Fields",
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test", "").
					WillReturnError(fmt.Errorf("user insert error"))
			},
			username: "test",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := u.DoRegister(context.Background(), tt.username, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPG_DoLogin(t *testing.T) {
	db, mock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := NewUserPG(db)

	tests := []struct {
		name         string
		mock         func()
		username     string
		password     string
		wantUsername string
		wantPassword string
		wantErr      bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()
				// mock.ExpectExec("INSERT INTO users").
				// 	WithArgs("test", "test").
				// 	WillReturnResult(sqlmock.NewResult(1, 1))
				db.Exec("INSERT INTO users VALUES ('test', 'test', 0, 0)")
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(1)
				mock.ExpectQuery("SELECT COUNT(*) FROM users").
					WithArgs("test", "test").WillReturnRows(rows)
				mock.ExpectCommit()
			},
			username:     "test",
			password:     "test",
			wantUsername: "test",
			wantPassword: "test",
		},
		// {
		// 	name: "Not Found",
		// 	mock: func() {
		// 		rows := sqlmock.NewRows([]string{"username", "password"}).
		// 			AddRow("test", "test")
		// 		mock.ExpectQuery("SELECT (.+) FROM users").
		// 			WithArgs("not", "found").WillReturnRows(rows)
		// 	},
		// 	username: "not",
		// 	password: "found",
		// 	wantErr:  true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := u.DoLogin(context.Background(), tt.username, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())

		})
	}
}
