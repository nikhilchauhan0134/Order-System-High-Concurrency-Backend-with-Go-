package db

import (
	"OrderSystemHighConcurrency/order-processor/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"database/sql"
	"strings"
)

// orderRepository implements contracts.Repository
type orderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB) contracts.Repository {
	return &orderRepository{db: db}
}

// SaveBatch inserts orders in bulk (SQL Server compatible)
func (r *orderRepository) SaveBatch(ctx context.Context, orders []*models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	var (
		query strings.Builder
		args  []interface{}
	)

	query.WriteString(`
		INSERT INTO orders (
			order_id, user_id, amount, currency, status,
			source, retry_count, created_at, updated_at
		) VALUES 
	`)

	for _, o := range orders {
		query.WriteString("(?,?,?,?,?,?,?,?,?),")

		args = append(args,
			o.OrderID,
			o.UserID,
			o.Amount,
			o.Currency,
			o.Status,
			o.Source,
			o.RetryCount,
			o.CreatedAt,
			o.UpdatedAt,
		)
	}

	// Remove trailing comma
	sqlQuery := strings.TrimSuffix(query.String(), ",")

	_, err := r.db.ExecContext(ctx, sqlQuery, args...)
	return err
}
