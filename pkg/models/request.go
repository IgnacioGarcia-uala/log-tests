package models

type Request struct {
	AccountId string  `json:"accountId"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
}
