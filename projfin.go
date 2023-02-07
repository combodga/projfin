package projfin

import "context"

var (
	Context context.Context
)

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type Accrual struct {
	OrderNum string  `json:"order"`
	Status   string  `json:"status"`
	Accrual  float64 `json:"accrual"`
}

type Order struct {
	OrderNumber string  `db:"order_number"`
	Username    string  `db:"username"`
	Status      string  `db:"status"`
	Accrual     float64 `db:"accrual"`
	UploadedAt  string  `db:"uploaded_at"`
}

type User struct {
	Username  string  `db:"username"`
	Password  string  `db:"password"`
	Balance   float64 `db:"balance"`
	Withdrawn float64 `db:"withdrawn"`
}

type Withdraw struct {
	OrderNumber string  `db:"order_number"`
	Username    string  `db:"username"`
	Sum         float64 `db:"sum"`
	ProcessedAt string  `db:"processed_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawShort struct {
	OrderNum string  `json:"order"`
	Sum      float64 `json:"sum"`
}

type WithdrawalsList struct {
	OrderNum    string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
