package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestWithdrawPG_ListWithdrawals(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	w := NewWithdrawPG(db)

	tests := []struct {
		name     string
		mock     func()
		username string
		wantErr  bool
	}{
		{
			name: "Ok",
			mock: func() {
				db.Exec("INSERT INTO withdrawals VALUES ('1234567', 'test', 0, NOW())")
				rows := sqlmock.NewRows([]string{"order_number", "username", "sum", "processed_at"}).AddRow(1)
				mock.ExpectQuery("SELECT * FROM withdrawals").WithArgs("1234567", "test").WillReturnRows(rows)
			},
			username: "test",
			wantErr:  false,
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

			_, err := w.ListWithdrawals(tt.username)
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
