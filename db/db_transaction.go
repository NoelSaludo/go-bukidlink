package db

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func QueryOrder(id string) (Order, error) {
	var order Order

	query := `
	SELECT user_id, item_id, amount, status, created_date 
	FROM "Order" 
	WHERE id = $1
	`

	row := db.QueryRow(query, id)
	if err := row.Err(); err != nil {
		return order, err
	}

	err := row.Scan(
		&order.UserId,
		&order.ItemId,
		&order.Amount,
		&order.Status,
		&order.CreatedDate,
	)

	order.Id = id

	if err != nil {
		return order, err
	}

	return order, nil
}

func InsertOrder(order Order) error {
	var err error

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO "Order"
	(id, user_id, item_id, amount, status, created_date)	
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec(order.Id, order.UserId, order.ItemId, order.Amount,
		order.Status, order.CreatedDate)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func DeleteOrder(id string) error {
	var err error

	query := `
	DELETE FROM "Order"	
	WHERE id = $1
	`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return err
}

func MoveOrderStatus(id, src, dst string) {
	var err error
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction: ", err)
	}

	if src == dst {
		log.Fatal("SRC and DST Columns are similar")
	}

	query := fmt.Sprintf(`DELETE FROM "%s" WHERE order_id = $1`, src)
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal("Delete query preparation failed: ", err)
	}
	_, err = tx.Stmt(stmt).Exec(id)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			log.Fatal("Rollback Failed: ", rberr)
		}
		log.Fatal("Failed to delete data: ", err)
	}

	query = fmt.Sprintf(`INSERT INTO "%s" (id, order_id) VALUES ($1, $2)`,
		dst)
	stmt, err = tx.Prepare(query)
	if err != nil {
		log.Fatal("Insert query failed: ", err)
	}
	_, err = tx.Stmt(stmt).Exec(uuid.New().String(), id)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			log.Fatal("Rollback Failed: ", rberr)
		}
		log.Fatal("Failed to delete data: ", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Commit failed: ", err)
	}
}
