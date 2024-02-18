package controller

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/forms"
	"github.com/ernestngugi/sil-backend/internal/model"
	"github.com/ernestngugi/sil-backend/internal/repos"
)

type (
	CustomerController interface {
		CustomerByName(ctx context.Context, dB db.DB, name string) (*model.Customer, error)
		CreateCustomer(ctx context.Context, dB db.DB, form *forms.CustomerCreateForm) (*model.Customer, error)
	}

	customerController struct {
		customerRepository repos.CustomerRepository
	}
)

func NewCustomerRepository(customerRepository repos.CustomerRepository) CustomerController {
	return &customerController{
		customerRepository: customerRepository,
	}
}

func (c *customerController) CustomerByName(
	ctx context.Context,
	dB db.DB,
	name string,
) (*model.Customer, error) {
	return c.customerRepository.CustomerByName(ctx, dB, strings.ToLower(name))
}

func (c *customerController) CreateCustomer(
	ctx context.Context,
	dB db.DB,
	form *forms.CustomerCreateForm,
) (*model.Customer, error) {

	_, err := c.customerRepository.CustomerByName(ctx, dB, strings.ToLower(form.Name))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			newCustomer := &model.Customer{
				Name: form.Name,
			}

			err := c.customerRepository.Save(ctx, dB, newCustomer)
			if err != nil {
				return &model.Customer{}, err
			}

			return newCustomer, nil
		}
		return &model.Customer{}, err
	}

	return &model.Customer{}, errors.New("customer exists")
}
