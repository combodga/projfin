package projfin

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
