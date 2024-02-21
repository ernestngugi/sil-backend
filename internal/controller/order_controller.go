package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/forms"
	"github.com/ernestngugi/sil-backend/internal/model"
	"github.com/ernestngugi/sil-backend/internal/repos"
	"github.com/ernestngugi/sil-backend/providers"
)

type (
	OrderController interface {
		CreateOrder(ctx context.Context, dB db.DB, form *forms.CreateOrderForm) (*model.Order, error)
		OrderByID(ctx context.Context, dB db.DB, orderID int64) (*model.Order, error)
	}

	orderController struct {
		atProvider         providers.ATProvider
		customerRepository repos.CustomerRepository
		orderRepository    repos.OrderRepository
	}
)

func NewOrderController(
	atProvider providers.ATProvider,
	customerRepository repos.CustomerRepository,
	orderRepository repos.OrderRepository,
) OrderController {
	return &orderController{
		atProvider:         atProvider,
		customerRepository: customerRepository,
		orderRepository:    orderRepository,
	}
}

func NewTestOrderController(
	atProvider providers.ATProvider,
) *orderController {
	return &orderController{
		atProvider:         atProvider,
		customerRepository: repos.NewCustomerRepository(),
		orderRepository:    repos.NewOrderRepository(),
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

	atRequest := &model.ATRequest{
		Number:  "254728389583",
		Message: fmt.Sprintf("your order %v has been received and is on your way", order.ID),
	}

	go c.sendOrderSMS(atRequest)

	return order, nil
}

func (c *orderController) sendOrderSMS(request *model.ATRequest) error {
	return c.atProvider.Send(request)
}
