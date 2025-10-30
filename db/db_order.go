package db

import (
	"database/sql"

	"github.com/google/uuid"
)

// QueryUsersOrders retrieves all orders for a given user, including their associated items.
func QueryUsersOrders(userId string) ([]Order, error) {
	var orders []Order

	query := `
	SELECT id, user_id, order_date, total_price, status
	FROM "Order"
	WHERE user_id = $1
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.Id, &order.UserId, &order.OrderDate, &order.TotalPrice, &order.Status)
		if err != nil {
			return nil, err
		}

		// For each order, get its items
		itemQuery := `
		SELECT id, order_id, item_id, quantity, price_at_purchase
		FROM "OrderItem"
		WHERE order_id = $1
		`
		itemRows, err := db.Query(itemQuery, order.Id)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item OrderItem
			err := itemRows.Scan(&item.Id, &item.OrderId, &item.ItemId, &item.Quantity, &item.PriceAtPurchase)
			if err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// InsertOrder creates a new order and its associated items in a single transaction.
func InsertOrder(order Order) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	// Insert into Order table
	orderId := uuid.New().String()
	orderQuery := `
	INSERT INTO "Order" (id, user_id, order_date, total_price, status)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(orderQuery, orderId, order.UserId, order.OrderDate, order.TotalPrice, order.Status)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Insert into OrderItem table
	itemStmt, err := tx.Prepare(`
	INSERT INTO "OrderItem" (id, order_id, item_id, quantity, price_at_purchase)
	VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	defer itemStmt.Close()

	for _, item := range order.Items {
		_, err := itemStmt.Exec(uuid.New().String(), orderId, item.ItemId, item.Quantity, item.PriceAtPurchase)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	return orderId, tx.Commit()
}

// UpdateOrderStatus updates the status of an existing order.
func UpdateOrderStatus(id string, status string) error {
	query := `
	UPDATE "Order"
	SET status = $1
	WHERE id = $2
	`

	result, err := db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
