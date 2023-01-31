package store

import (
	"context"
	"testing"

	"github.com/combodga/projfin"
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
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("test", "test")
			},
			username: "test",
			password: "test",
		},
		// {
		// 	name: "Empty Fields",
		// 	mock: func() {
		// 		mock.ExpectQuery("INSERT INTO users").
		// 			WithArgs("test", "")
		// 	},
		// 	username: "test",
		// 	password: "",
		// 	wantErr:  true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			projfin.Context = context.Background()
			err := u.DoRegister(tt.username, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

/*
func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type args struct {
		username string
		password string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    todo.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"}).
					AddRow(1, "Test", "test", "password")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("test", "password").WillReturnRows(rows)
			},
			input: args{"test", "password"},
			want: todo.User{
				Id:       1,
				Name:     "Test",
				Username: "test",
				Password: "password",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("not", "found").WillReturnRows(rows)
			},
			input:   args{"not", "found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUser(tt.input.username, tt.input.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
*/
