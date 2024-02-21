package model

import "time"

type Order struct {
	ID          int64     `json:"id"`
	Amount      float64   `json:"amount"`
	CustomerID  int64     `json:"customer_id"`
	Item        string    `json:"item"`
	DateCreated time.Time `json:"date_created"`
}

type ATRequest struct {
	Number  string
	Message string
}
