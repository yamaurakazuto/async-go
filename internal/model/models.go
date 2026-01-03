package model

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID int64 `json:"user_id"`
}

type BalanceResponse struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

type TransferRequest struct {
	FromUserID      int64   `json:"from_user_id"`
	ToAccountNumber string  `json:"to_account_number"`
	Amount          float64 `json:"amount"`
}

type TransferResponse struct {
	TransferID int64 `json:"transfer_id"`
}
