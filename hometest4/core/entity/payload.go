package entity

// register user payload
type RegisterPayload struct {
	Username string `json:"username" validate:"required,max=30,username"`
}

// transer payload
type TransferPayload struct {
	ToUsername string  `json:"to_username" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
}

// balance topup payload
type BalanceTopupPayload struct {
	Amount float64 `json:"amount" validate:"required,max=9999999,min=1"`
}
