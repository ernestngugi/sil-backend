package controller

import (
	"context"
	"testing"

	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/forms"
	"github.com/ernestngugi/sil-backend/internal/model"
	"github.com/ernestngugi/sil-backend/internal/repos"
	"github.com/stretchr/testify/assert"
)

func TestOrderController(t *testing.T) {

	t.Setenv("DATABASE_URL", "postgres://savannah:password@localhost:5433/savannah?sslmode=disable")

	ctx := context.Background()

	dB := db.InitDB()
	defer dB.Close()

	customerRepository := repos.NewCustomerRepository()
	orderRepository := repos.NewOrderRepository()

	orderController := NewOrderController(customerRepository, orderRepository)

	t.Run("cannot create an order if credentials are missing", func(t *testing.T) {

		form := &forms.CreateOrderForm{
			Amount: 100,
			Item:   "item",
		}

		_, err := orderController.CreateOrder(ctx, dB, form)
		assert.Error(t, err)

		assert.Equal(t, "customer credential missing", err.Error())

		clearOrderTable(ctx, dB)
	})

	t.Run("can create an order", func(t *testing.T) {

		customer := model.BuildCustomer()

		err := customerRepository.Save(ctx, dB, customer)
		assert.NoError(t, err)

		ctx = context.WithValue(ctx, model.CustomerKeyName, customer.Name)

		form := &forms.CreateOrderForm{
			Amount: 100,
			Item:   "item",
		}

		order, err := orderController.CreateOrder(ctx, dB, form)
		assert.NoError(t, err)
		assert.NotZero(t, order.ID)
		assert.NotZero(t, order.DateCreated)
		assert.Equal(t, order.CustomerID, customer.ID)
		assert.Equal(t, order.Amount, 100.00)
		assert.Equal(t, order.Item, "item")

		clearOrderTable(ctx, dB)
		clearCustomerTable(ctx, dB)
	})
}

func clearOrderTable(ctx context.Context, dB db.DB) {
	dB.ExecContext(ctx, "DELETE FROM orders")
	dB.ExecContext(ctx, "ALTER SEQUENCE orders_id_seq RESTART WITH 1")
}
