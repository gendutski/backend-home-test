package entity

// register user ok response
type RegisterResponse struct {
	Token string `json:"token"`
}

// balance read ok response
type BalanceReadResponse struct {
	Balance float64 `json:"balance"`
}

// top transaction by user amount size
type TopTransactionResponse struct {
	Username         string  `json:"username"`
	TransactionValue float64 `json:"amount"`
}

// top transaction by sum of debit amount
type OverallTopTransactionResponse struct {
	Username         string  `json:"username"`
	TransactionValue float64 `json:"transacted_value"`
}
