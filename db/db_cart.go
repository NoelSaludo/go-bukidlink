package db

import (
	"database/sql"

	"github.com/google/uuid"
)

// GetCartByUserID retrieves the cart for a given user.
// If the user does not have a cart, one is created.
func GetCartByUserID(userId string) (Cart, error) {
	var cart Cart

	query := `SELECT id, user_id, grand_total, created_at FROM "Cart" WHERE user_id = $1`
	row := db.QueryRow(query, userId)

	err := row.Scan(&cart.Id, &cart.UserId, &cart.GrandTotal, &cart.CreatedAt)
	if err == sql.ErrNoRows {
		// No cart found, create one
		newCartId := uuid.New().String()
		insertQuery := `INSERT INTO "Cart" (id, user_id, grand_total) VALUES ($1, $2, 0.0)`
		_, err := db.Exec(insertQuery, newCartId, userId)
		if err != nil {
			return Cart{}, err
		}
		// Return the newly created, empty cart
		return GetCartByUserID(userId)
	} else if err != nil {
		return Cart{}, err
	}

	// Cart found, now get its items
	itemQuery := `SELECT id, cart_id, item_id, quantity FROM "CartItem" WHERE cart_id = $1`
	itemRows, err := db.Query(itemQuery, cart.Id)
	if err != nil {
		return cart, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item CartItem
		err := itemRows.Scan(&item.Id, &item.CartId, &item.ItemId, &item.Quantity)
		if err != nil {
			return cart, err
		}
		cart.Items = append(cart.Items, item)
	}

	return cart, nil
}

// AddItemToCart adds an item to a cart. If the item already exists, it updates the quantity.
func AddItemToCart(cartId string, itemId string, quantity int) error {
	var existingQuantity int
	query := `SELECT quantity FROM "CartItem" WHERE cart_id = $1 AND item_id = $2`
	row := db.QueryRow(query, cartId, itemId)
	err := row.Scan(&existingQuantity)

	if err == sql.ErrNoRows {
		// Item does not exist, insert it
		insertQuery := `INSERT INTO "CartItem" (id, cart_id, item_id, quantity) VALUES ($1, $2, $3, $4)`
		_, err := db.Exec(insertQuery, uuid.New().String(), cartId, itemId, quantity)
		return err
	} else if err != nil {
		return err
	}

	// Item exists, update the quantity
	newQuantity := existingQuantity + quantity
	updateQuery := `UPDATE "CartItem" SET quantity = $1 WHERE cart_id = $2 AND item_id = $3`
	_, err = db.Exec(updateQuery, newQuantity, cartId, itemId)
	return err
}

// RemoveItemFromCart removes an item from a cart completely.
func RemoveItemFromCart(cartItemId string) error {
	query := `DELETE FROM "CartItem" WHERE id = $1`
	_, err := db.Exec(query, cartItemId)
	return err
}

// UpdateCartItemQuantity updates the quantity for a cart item. If quantity <= 0,
// the cart item is removed.
func UpdateCartItemQuantity(cartItemId string, quantity int) error {
	if quantity <= 0 {
		// remove the item
		return RemoveItemFromCart(cartItemId)
	}

	query := `UPDATE "CartItem" SET quantity = $1 WHERE id = $2`
	_, err := db.Exec(query, quantity, cartItemId)
	return err
}
