package controller

import (
	"context"
	"testing"

	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/forms"
	"github.com/ernestngugi/sil-backend/internal/repos"
	"github.com/stretchr/testify/assert"
)

func TestCustomerController(t *testing.T) {

	t.Setenv("DATABASE_URL", "postgres://savannah:password@localhost:5433/savannah?sslmode=disable")

	ctx := context.Background()

	dB := db.InitDB()
	defer dB.Close()

	customerRepository := repos.NewCustomerRepository()

	customerController := NewCustomerController(customerRepository)

	t.Run("can create a customer", func(t *testing.T) {

		form := &forms.CustomerCreateForm{
			Name: "test",
		}

		customer, err := customerController.CreateCustomer(ctx, dB, form)
		assert.NoError(t, err)

		assert.NotZero(t, customer.ID)
		assert.Equal(t, "test", customer.Name)
		assert.NotZero(t, customer.DateCreated)
		assert.NotZero(t, customer.DateModified)

		clearCustomerTable(ctx, dB)
	})

	t.Run("cannot create customers with similar name", func(t *testing.T) {

		form := &forms.CustomerCreateForm{
			Name: "test",
		}

		customer, err := customerController.CreateCustomer(ctx, dB, form)
		assert.NoError(t, err)

		assert.NotZero(t, customer.ID)
		assert.Equal(t, "test", customer.Name)
		assert.NotZero(t, customer.DateCreated)
		assert.NotZero(t, customer.DateModified)

		_, err = customerController.CreateCustomer(ctx, dB, form)
		assert.Error(t, err)

		assert.Equal(t, "customer exists", err.Error())

		clearCustomerTable(ctx, dB)
	})

	t.Run("can get customer by name", func(t *testing.T) {

		form := &forms.CustomerCreateForm{
			Name: "test",
		}

		customer, err := customerController.CreateCustomer(ctx, dB, form)
		assert.NoError(t, err)

		assert.NotZero(t, customer.ID)
		assert.Equal(t, "test", customer.Name)
		assert.NotZero(t, customer.DateCreated)
		assert.NotZero(t, customer.DateModified)

		foundCustomer, err := customerRepository.CustomerByName(ctx, dB, "test")
		assert.NoError(t, err)
		assert.Equal(t, customer.ID, foundCustomer.ID)
		assert.Equal(t, "test", foundCustomer.Name)

		clearCustomerTable(ctx, dB)
	})
}

func clearCustomerTable(ctx context.Context, dB db.DB) {
	dB.ExecContext(ctx, "DELETE FROM customers")
	dB.ExecContext(ctx, "ALTER SEQUENCE customers_id_seq RESTART WITH 1")
}
