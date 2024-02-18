package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/forms"
	"github.com/ernestngugi/sil-backend/internal/model"
	"github.com/ernestngugi/sil-backend/internal/repos"
)

type (
	OrderController interface {
		CreateOrder(ctx context.Context, dB db.DB, form *forms.CreateOrderForm) (*model.Order, error)
		OrderByID(ctx context.Context, dB db.DB, orderID int64) (*model.Order, error)
	}

	orderController struct {
		customerRepository repos.CustomerRepository
		orderRepository    repos.OrderRepository
	}
)

func NewOrderController(
	customerRepository repos.CustomerRepository,
	orderRepository repos.OrderRepository,
) OrderController {
	return &orderController{
		customerRepository: customerRepository,
		orderRepository:    orderRepository,
	}
}

func (c *orderController) OrderByID(
	ctx context.Context,
	dB db.DB,
	orderID int64,
) (*model.Order, error) {
	return c.orderRepository.OrderByID(ctx, dB, orderID)
}

func (c *orderController) CreateOrder(
	ctx context.Context,
	dB db.DB,
	form *forms.CreateOrderForm,
) (*model.Order, error) {

	exist := ctx.Value(model.CustomerKeyName)
	if exist == nil {
		return &model.Order{}, errors.New("customer credential missing")
	}

	name, ok := exist.(string)
	if !ok {
		return &model.Order{}, errors.New("customer credential missing")
	}

	customer, err := c.customerRepository.CustomerByName(ctx, dB, strings.ToLower(name))
	if err != nil {
		return &model.Order{}, err
	}

	order := &model.Order{
		CustomerID: customer.ID,
		Amount:     form.Amount,
		Item:       form.Item,
	}

	err = c.orderRepository.Save(ctx, dB, order)
	if err != nil {
		return &model.Order{}, err
	}

	return order, nil
}
