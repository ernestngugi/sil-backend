package model

import (
	"time"

	"syreclabs.com/go/faker"
)

type Customer struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified"`
}

type CustomerKey string

const CustomerKeyName CustomerKey = "customer_name"

func BuildCustomer() *Customer {
	return &Customer{
		Name: faker.Internet().Email(),
	}
}
