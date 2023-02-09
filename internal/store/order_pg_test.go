package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestOrderPG_MakeOrder(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	o := NewOrderPG(db)

	tests := []struct {
		name        string
		mock        func()
		orderNumber string
		username    string
		wantErr     bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("INSERT INTO orders").
					WithArgs("1234567", "test").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			orderNumber: "1234567",
			username:    "test",
		},
		{
			name: "Empty Fields",
			mock: func() {
				mock.ExpectExec("INSERT INTO orders").
					WithArgs("", "").
					WillReturnError(fmt.Errorf("order insert error"))
			},
			orderNumber: "",
			username:    "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := o.MakeOrder(context.Background(), tt.username, tt.orderNumber)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderPG_CheckOrder(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	o := NewOrderPG(db)

	tests := []struct {
		name            string
		mock            func()
		orderNumber     string
		username        string
		wantOrderNumber string
		wantUsername    string
		wantErr         bool
	}{
		{
			name: "Ok",
			mock: func() {
				db.Exec("INSERT INTO orders VALUES ('1234567', 'test', 'NEW', 0, NOW())")
				rows := sqlmock.NewRows([]string{"order_number", "username", "status", "accrual", "uploaded_at"}).AddRow(1)
				mock.ExpectQuery("SELECT * FROM orders").WithArgs("1234567", "test").WillReturnRows(rows)
			},
			orderNumber:     "",
			username:        "test",
			wantOrderNumber: "",
			wantUsername:    "test",
			wantErr:         false,
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

			_, err := o.CheckOrder(tt.username, tt.orderNumber)
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
