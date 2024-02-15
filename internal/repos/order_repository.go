package repos

import (
	"context"
	"time"

	"github.com/ernestngugi/sil-backend/internal/db"
	"github.com/ernestngugi/sil-backend/internal/model"
)

const (
	insertOrderSQL  = "INSERT INTO orders(item, amount, customer_id, date_created) VALUES ($1, $2, $3, $4) RETURNING id"
	selectOrderSQL  = "SELECT id, item, amount, customer_id, date_created FROM orders"
	getOrderByIDSQL = selectOrderSQL + " WHERE id = $1"
)

type (
	OrderRepository interface {
		OrderByID(ctx context.Context, operations db.SQLOperations, orderID int64) (*model.Order, error)
		Save(ctx context.Context, operations db.SQLOperations, order *model.Order) error
	}

	orderRepository struct{}
)

func NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

func (r *orderRepository) OrderByID(
	ctx context.Context,
	operations db.SQLOperations,
	orderID int64,
) (*model.Order, error) {

	row := operations.QueryRowContext(ctx, getOrderByIDSQL, orderID)

	var order model.Order

	err := row.Scan(&order.ID, &order.Item, &order.Amount, &order.CustomerID, &order.DateCreated)
	if err != nil {
		return &model.Order{}, err
	}

	return &order, nil
}

func (r *orderRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	order *model.Order,
) error {

	timeNow := time.Now()
	order.DateCreated = timeNow

	err := operations.QueryRowContext(
		ctx,
		insertOrderSQL,
		order.Item,
		order.Amount,
		order.CustomerID,
		order.DateCreated,
	).Scan(&order.ID)
	if err != nil {
		return err
	}

	return nil
}
