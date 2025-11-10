package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	HOST     string = os.Getenv("DBHOST")
	PORT     string = os.Getenv("DBPORT")
	USER     string = os.Getenv("DBUSER")
	PASSWORD string = os.Getenv("DBPASSWORD")
	DATABASE string = os.Getenv("DATABASE")
)

var db *sql.DB

func SetupDatabase() error {
	var err error

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST,
		PORT,
		USER,
		PASSWORD,
		DATABASE,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return nil
}

func Ping() error {
	return db.Ping()
}

// must be used only with items and with proper order
func getItemsFromRow(rows *sql.Rows, items []Item) []Item {
	for rows.Next() {
		var item Item

		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.CostPKilo,
			&item.Category,
			&item.Amount,
			&item.ImgPath,
		)

		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	return items
}

func getReviewFromRow(rows *sql.Rows, reviews []Review) ([]Review, error) {
	for rows.Next() {
		var review Review

		err := rows.Scan(
			&review.Id,
			&review.UserId,
			&review.ItemId,
			&review.Content,
			&review.Rating)

		if err != nil {
			return reviews, err
		}

		reviews = append(reviews, review)
	}
	return reviews, nil
}

// rollbackAndReturn attempts to rollback the provided transaction and
// returns the rollback error if it occurs; otherwise it returns the
// original error that triggered the rollback.
//
// Usage:
//
//	if err := someDBOp(); err != nil {
//	    return rollbackAndReturn(tx, err)
//	}
func rollbackAndReturn(tx *sql.Tx, origErr error) error {
	if rberr := tx.Rollback(); rberr != nil {
		return rberr
	}
	return origErr
}

func GetDB() *sql.DB {
	return db
}
