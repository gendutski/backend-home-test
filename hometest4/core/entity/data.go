package entity

import "time"

type TrxFlow int

const (
	TrxFlowDebit TrxFlow = iota
	TrxFlowCredit
)

// user data
type User struct {
	Username string
	Balance  float64
}

// transaction data
type Transaction struct {
	Username string
	Amount   float64
	Flow     TrxFlow
	TrxDate  time.Time
}
