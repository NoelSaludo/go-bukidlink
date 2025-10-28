package db

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func QueryTransaction(id string) (Transaction, error) {
	var transaction Transaction

	query := `
	SELECT user_id, item_id, amount, created_date 
	FROM "Transaction" 
	WHERE id = $1
	`

	row := db.QueryRow(query, id)
	if err := row.Err(); err != nil {
		fmt.Errorf("Error querying row: %s", err.Error())
		return transaction, err
	}

	err := row.Scan(
		&transaction.UserId,
		&transaction.ItemId,
		&transaction.amount,
		&transaction.CreatedDate,
	)

	transaction.Id = id

	if err != nil {
		fmt.Errorf("Error scanning: %s", err.Error())
		return transaction, err
	}

	return transaction, nil
}

func InsertTransaction(trans Transaction) error {
	var err error

	query := `
	INSERT INTO "Transaction"
	(id, user_id, item_id, amount, created_date)	
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err = db.Exec(query, trans.Id, trans.UserId, trans.ItemId, trans.amount,
		trans.CreatedDate)
	if err != nil {
		fmt.Errorf("Inserting error: %s", err.Error())
		return err
	}

	return nil
}

func DeleteTransaction(id string) error {
	var err error

	query := `
	DELETE FROM "Transaction"	
	WHERE id = $1
	`

	_, err = db.Exec(query, id)
	if err != nil {
		fmt.Errorf("Deletion Failed: %s", err.Error())
		return err
	}

	return err
}

func MoveTransactionStatus(id, src, dst string) {
	var err error
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction: ", err)
	}

	if src == dst {
		log.Fatal("SRC and DST Columns are similar")
	}

	query := fmt.Sprintf(`DELETE FROM "%s" WHERE transaction_id = $1`, src)
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

	query = fmt.Sprintf(`INSERT INTO "%s" (id, transaction_id) VALUES ($1, $2)`,
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
