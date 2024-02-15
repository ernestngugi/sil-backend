package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/model"
)

const (
	insertCustomerSQL  = "INSERT INTO customers (name, date_created, date_modified) VALUES ($1, $2, $3) RETURNING id"
	selectCustomerSQL  = "SELECT id, name, date_created, date_modified FROM customers"
	getCustomerByIDSQL = selectCustomerSQL + " WHERE id = $1"
)

type (
	CustomerRepository interface {
		CustomerByID(ctx context.Context, operations db.SQLOperations, customerID int64) (*model.Customer, error)
		Save(ctx context.Context, operations db.SQLOperations, customer *model.Customer) error
	}

	customerRepository struct{}
)

func NewCustomerRepository() CustomerRepository {
	return &customerRepository{}
}

func (r *customerRepository) CustomerByID(
	ctx context.Context,
	operations db.SQLOperations,
	customerID int64,
) (*model.Customer, error) {

	row := operations.QueryRowContext(ctx, getCustomerByIDSQL, customerID)

	var customer model.Customer

	err := row.Scan(&customer.ID, &customer.Name, &customer.DateCreated, &customer.DateModified)
	if err != nil {
		return &model.Customer{}, err
	}

	return &customer, nil
}

func (r *customerRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	customer *model.Customer,
) error {

	timeNow := time.Now()
	customer.DateCreated = timeNow
	customer.DateModified = timeNow

	if customer.ID == 0 {

		err := operations.QueryRowContext(
			ctx,
			insertCustomerSQL,
			customer.Name,
			customer.DateCreated,
			customer.DateModified,
		).Scan(&customer.ID)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("update customer forbidden")
}
