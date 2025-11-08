package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GetUserBalance retrieves the balance for a specific user.
func GetUserBalance(userID string) (UserBalance, error) {
	var balance UserBalance
	query := `SELECT id, user_id, balance, currency, updated_at, created_at 
	          FROM "UserBalance" WHERE user_id = $1`

	err := db.QueryRow(query, userID).Scan(
		&balance.ID,
		&balance.UserID,
		&balance.Balance,
		&balance.Currency,
		&balance.UpdatedAt,
		&balance.CreatedAt,
	)

	return balance, err
}

// CreateUserBalance initializes a balance record for a new user.
func CreateUserBalance(userID string) error {
	balanceID := uuid.New().String()
	query := `INSERT INTO "UserBalance" (id, user_id, balance, currency, created_at, updated_at)
	          VALUES ($1, $2, 0.00, 'PHP', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := db.Exec(query, balanceID, userID)
	return err
}

// GetPaymentTransaction retrieves a specific transaction by ID.
func GetPaymentTransaction(transactionID string) (PaymentTransaction, error) {
	var txn PaymentTransaction
	query := `SELECT id, user_id, order_id, transaction_type, amount, balance_before, 
	          balance_after, status, payment_method, reference_number, description, 
	          created_at, updated_at, completed_at
	          FROM "PaymentTransaction" WHERE id = $1`

	err := db.QueryRow(query, transactionID).Scan(
		&txn.ID,
		&txn.UserID,
		&txn.OrderID,
		&txn.TransactionType,
		&txn.Amount,
		&txn.BalanceBefore,
		&txn.BalanceAfter,
		&txn.Status,
		&txn.PaymentMethod,
		&txn.ReferenceNumber,
		&txn.Description,
		&txn.CreatedAt,
		&txn.UpdatedAt,
		&txn.CompletedAt,
	)

	return txn, err
}

// GetUserTransactions retrieves all transactions for a user, ordered by creation date.
func GetUserTransactions(userID string, limit, offset int) ([]PaymentTransaction, error) {
	query := `SELECT id, user_id, order_id, transaction_type, amount, balance_before, 
	          balance_after, status, payment_method, reference_number, description, 
	          created_at, updated_at, completed_at
	          FROM "PaymentTransaction" 
	          WHERE user_id = $1 
	          ORDER BY created_at DESC 
	          LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []PaymentTransaction
	for rows.Next() {
		var txn PaymentTransaction
		err := rows.Scan(
			&txn.ID,
			&txn.UserID,
			&txn.OrderID,
			&txn.TransactionType,
			&txn.Amount,
			&txn.BalanceBefore,
			&txn.BalanceAfter,
			&txn.Status,
			&txn.PaymentMethod,
			&txn.ReferenceNumber,
			&txn.Description,
			&txn.CreatedAt,
			&txn.UpdatedAt,
			&txn.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}

	return transactions, nil
}

// ProcessDeposit adds funds to a user's balance.
func ProcessDeposit(userID string, amount float64, paymentMethod, referenceNumber, description string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Lock and get current balance
	var balance UserBalance
	err = tx.QueryRow(`SELECT id, user_id, balance FROM "UserBalance" WHERE user_id = $1 FOR UPDATE`,
		userID).Scan(&balance.ID, &balance.UserID, &balance.Balance)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	balanceBefore := balance.Balance
	balanceAfter := balanceBefore + amount
	transactionID := uuid.New().String()
	now := time.Now()

	// Create transaction record
	_, err = tx.Exec(`INSERT INTO "PaymentTransaction" 
		(id, user_id, transaction_type, amount, balance_before, balance_after, status, 
		 payment_method, reference_number, description, created_at, updated_at, completed_at)
		VALUES ($1, $2, 'deposit', $3, $4, $5, 'completed', $6, $7, $8, $9, $9, $9)`,
		transactionID, userID, amount, balanceBefore, balanceAfter,
		paymentMethod, referenceNumber, description, now)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	// Update balance
	_, err = tx.Exec(`UPDATE "UserBalance" SET balance = $1, updated_at = $2 WHERE user_id = $3`,
		balanceAfter, now, userID)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	return tx.Commit()
}

// ProcessOrderPayment deducts funds from user balance for an order.
func ProcessOrderPayment(userID, orderID string, amount float64, description string) (PaymentTransaction, error) {
	var txn PaymentTransaction

	tx, err := db.Begin()
	if err != nil {
		return txn, err
	}

	// Ensure the referenced order exists in the database. If the order doesn't
	// exist at the time of inserting the payment transaction, the DB will raise
	// a foreign key violation. Check early and return a clear error.
	var existingOrderID string
	err = tx.QueryRow(`SELECT id FROM "Order" WHERE id = $1`, orderID).Scan(&existingOrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return txn, rollbackAndReturn(tx, fmt.Errorf("order not found: %s", orderID))
		}
		return txn, rollbackAndReturn(tx, err)
	}

	// Lock and get current balance
	var balance UserBalance
	err = tx.QueryRow(`SELECT id, user_id, balance FROM "UserBalance" WHERE user_id = $1 FOR UPDATE`,
		userID).Scan(&balance.ID, &balance.UserID, &balance.Balance)
	if err != nil {
		return txn, rollbackAndReturn(tx, err)
	}

	// Check sufficient balance
	if balance.Balance < amount {
		return txn, rollbackAndReturn(tx, sql.ErrNoRows) // Return standard error for insufficient funds
	}

	balanceBefore := balance.Balance
	balanceAfter := balanceBefore - amount
	transactionID := uuid.New().String()
	now := time.Now()

	// Create transaction record
	err = tx.QueryRow(`INSERT INTO "PaymentTransaction" 
		(id, user_id, order_id, transaction_type, amount, balance_before, balance_after, status, 
		 payment_method, description, created_at, updated_at, completed_at)
		VALUES ($1, $2, $3, 'order_payment', $4, $5, $6, 'completed', 'balance', $7, $8, $8, $8)
		RETURNING id, user_id, order_id, transaction_type, amount, balance_before, balance_after, 
		          status, payment_method, reference_number, description, created_at, updated_at, completed_at`,
		transactionID, userID, orderID, amount, balanceBefore, balanceAfter, description, now).Scan(
		&txn.ID,
		&txn.UserID,
		&txn.OrderID,
		&txn.TransactionType,
		&txn.Amount,
		&txn.BalanceBefore,
		&txn.BalanceAfter,
		&txn.Status,
		&txn.PaymentMethod,
		&txn.ReferenceNumber,
		&txn.Description,
		&txn.CreatedAt,
		&txn.UpdatedAt,
		&txn.CompletedAt,
	)
	if err != nil {
		return txn, rollbackAndReturn(tx, err)
	}

	// Update balance
	_, err = tx.Exec(`UPDATE "UserBalance" SET balance = $1, updated_at = $2 WHERE user_id = $3`,
		balanceAfter, now, userID)
	if err != nil {
		return txn, rollbackAndReturn(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return txn, err
	}

	return txn, nil
}

// ProcessRefund returns funds to user balance.
func ProcessRefund(userID, orderID string, amount float64, description string) (PaymentTransaction, error) {
	var txn PaymentTransaction

	tx, err := db.Begin()
	if err != nil {
		return txn, err
	}

	// Lock and get current balance
	var balance UserBalance
	err = tx.QueryRow(`SELECT id, user_id, balance FROM "UserBalance" WHERE user_id = $1 FOR UPDATE`,
		userID).Scan(&balance.ID, &balance.UserID, &balance.Balance)
	if err != nil {
		return txn, rollbackAndReturn(tx, err)
	}

	balanceBefore := balance.Balance
	balanceAfter := balanceBefore + amount
	transactionID := uuid.New().String()
	now := time.Now()

	// Create transaction record
	err = tx.QueryRow(`INSERT INTO "PaymentTransaction" 
		(id, user_id, order_id, transaction_type, amount, balance_before, balance_after, status, 
		 payment_method, description, created_at, updated_at, completed_at)
		VALUES ($1, $2, $3, 'refund', $4, $5, $6, 'completed', 'balance', $7, $8, $8, $8)
		RETURNING id, user_id, order_id, transaction_type, amount, balance_before, balance_after, 
		          status, payment_method, reference_number, description, created_at, updated_at, completed_at`,
		transactionID, userID, orderID, amount, balanceBefore, balanceAfter, description, now).Scan(
		&txn.ID,
		&txn.UserID,
		&txn.OrderID,
		&txn.TransactionType,
		&txn.Amount,
		&txn.BalanceBefore,
		&txn.BalanceAfter,
		&txn.Status,
		&txn.PaymentMethod,
		&txn.ReferenceNumber,
		&txn.Description,
		&txn.CreatedAt,
		&txn.UpdatedAt,
		&txn.CompletedAt,
	)
	if err != nil {
		return txn, rollbackAndReturn(tx, err)
	}

	// Update balance
	_, err = tx.Exec(`UPDATE "UserBalance" SET balance = $1, updated_at = $2 WHERE user_id = $3`,
		balanceAfter, now, userID)
	if err != nil {
		return txn, rollbackAndReturn(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return txn, err
	}

	return txn, nil
}

// ProcessWithdrawal removes funds from user balance.
func ProcessWithdrawal(userID string, amount float64, paymentMethod, referenceNumber, description string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Lock and get current balance
	var balance UserBalance
	err = tx.QueryRow(`SELECT id, user_id, balance FROM "UserBalance" WHERE user_id = $1 FOR UPDATE`,
		userID).Scan(&balance.ID, &balance.UserID, &balance.Balance)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	// Check sufficient balance
	if balance.Balance < amount {
		return rollbackAndReturn(tx, sql.ErrNoRows)
	}

	balanceBefore := balance.Balance
	balanceAfter := balanceBefore - amount
	transactionID := uuid.New().String()
	now := time.Now()

	// Create transaction record
	_, err = tx.Exec(`INSERT INTO "PaymentTransaction" 
		(id, user_id, transaction_type, amount, balance_before, balance_after, status, 
		 payment_method, reference_number, description, created_at, updated_at, completed_at)
		VALUES ($1, $2, 'withdrawal', $3, $4, $5, 'completed', $6, $7, $8, $9, $9, $9)`,
		transactionID, userID, amount, balanceBefore, balanceAfter,
		paymentMethod, referenceNumber, description, now)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	// Update balance
	_, err = tx.Exec(`UPDATE "UserBalance" SET balance = $1, updated_at = $2 WHERE user_id = $3`,
		balanceAfter, now, userID)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	return tx.Commit()
}
