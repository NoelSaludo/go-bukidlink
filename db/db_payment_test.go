package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetUserBalance tests retrieving an existing user balance
func TestGetUserBalance(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	balance, err := GetUserBalance(userID)

	require.NoError(t, err)
	assert.Equal(t, userID, balance.UserID)
	assert.Equal(t, "PHP", balance.Currency)
	assert.GreaterOrEqual(t, balance.Balance, float64(0))
	assert.NotEmpty(t, balance.ID)
	assert.NotZero(t, balance.CreatedAt)
	assert.NotZero(t, balance.UpdatedAt)
}

// TestGetUserBalance_NotFound tests retrieving balance for non-existent user
func TestGetUserBalance_NotFound(t *testing.T) {
	_ = SetupDatabase()

	_, err := GetUserBalance("00000000-0000-0000-0000-000000000000")

	require.Error(t, err)
}

// TestCreateUserBalance tests creating a new balance record
func TestCreateUserBalance(t *testing.T) {
	_ = SetupDatabase()

	// Create a unique user ID for testing
	existingUserID := "9ae195a0-05ff-446b-99c0-e6f09a0150d1"

	// Create balance
	err := CreateUserBalance(existingUserID)
	require.NoError(t, err)

	// Verify the balance was created
	balance, err := GetUserBalance(existingUserID)
	require.NoError(t, err)
	assert.Equal(t, existingUserID, balance.UserID)
	assert.Equal(t, float64(0), balance.Balance)
	assert.Equal(t, "PHP", balance.Currency)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "UserBalance" WHERE user_id = $1`, existingUserID)
}

// TestProcessDeposit tests adding funds to a user balance
func TestProcessDeposit(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Get initial balance
	initialBalance, err := GetUserBalance(userID)
	require.NoError(t, err)

	depositAmount := 500.00
	paymentMethod := "gcash"
	referenceNumber := "GC" + uuid.New().String()[:8]
	description := "Test deposit"

	// Process deposit
	err = ProcessDeposit(userID, depositAmount, paymentMethod, referenceNumber, description)

	require.NoError(t, err)

	// Verify balance was updated
	newBalance, err := GetUserBalance(userID)
	require.NoError(t, err)
	assert.Equal(t, initialBalance.Balance+depositAmount, newBalance.Balance)

	// Cleanup - reverse the deposit
	_ = ProcessWithdrawal(userID, depositAmount, "test", "cleanup", "Test cleanup")
}

// TestProcessOrderPayment tests paying for an order
func TestProcessOrderPayment(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// First deposit funds to ensure sufficient balance
	depositAmount := 1000.00
	err := ProcessDeposit(userID, depositAmount, "test", "TEST"+uuid.New().String()[:8], "Setup funds")
	require.NoError(t, err)

	// Get balance after deposit
	balanceAfterDeposit, err := GetUserBalance(userID)
	require.NoError(t, err)

	// Use existing test order ID from init-test-db.sql (belongs to JohnDoe)
	orderID := "11111111-1111-1111-1111-111111111111"
	paymentAmount := 250.00
	description := "Payment for test order"

	// Process payment
	txn, err := ProcessOrderPayment(userID, orderID, paymentAmount, description)

	require.NoError(t, err)
	assert.NotEmpty(t, txn.ID)
	assert.Equal(t, userID, txn.UserID)
	assert.NotNil(t, txn.OrderID)
	assert.Equal(t, orderID, *txn.OrderID)
	assert.Equal(t, "order_payment", txn.TransactionType)
	assert.Equal(t, paymentAmount, txn.Amount)
	assert.Equal(t, balanceAfterDeposit.Balance, txn.BalanceBefore)
	assert.Equal(t, balanceAfterDeposit.Balance-paymentAmount, txn.BalanceAfter)
	assert.Equal(t, "completed", txn.Status)
	assert.Equal(t, "balance", txn.PaymentMethod)
	assert.NotZero(t, txn.CreatedAt)
	assert.NotNil(t, txn.CompletedAt)

	// Verify balance was updated
	newBalance, err := GetUserBalance(userID)
	require.NoError(t, err)
	assert.Equal(t, balanceAfterDeposit.Balance-paymentAmount, newBalance.Balance)

	// Cleanup - refund and withdraw
	_, _ = ProcessRefund(userID, orderID, paymentAmount, "Test cleanup")
	_ = ProcessWithdrawal(userID, depositAmount, "test", "cleanup", "Test cleanup")
}

// TestProcessOrderPayment_InsufficientFunds tests payment with insufficient balance
func TestProcessOrderPayment_InsufficientFunds(t *testing.T) {
	_ = SetupDatabase()

	// Create a new user with zero balance
	newUserID := "543255dd-5325-4d3f-bcd2-ee6f8ac87e2e"
	err := CreateUserBalance(newUserID)
	require.NoError(t, err)

	// Try to pay more than available balance
	orderID := uuid.New().String()
	_, err = ProcessOrderPayment(newUserID, orderID, 100.00, "Test payment")

	require.Error(t, err)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "UserBalance" WHERE user_id = $1`, newUserID)
}

// TestProcessRefund tests refunding funds to user balance
func TestProcessRefund(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Get initial balance
	initialBalance, err := GetUserBalance(userID)
	require.NoError(t, err)

	orderID := "983a8b8f-2fcf-45ce-b26f-5965ae5fef8b"
	refundAmount := 150.00
	description := "Test refund for cancelled order"

	// Process refund
	txn, err := ProcessRefund(userID, orderID, refundAmount, description)

	require.NoError(t, err)
	assert.NotEmpty(t, txn.ID)
	assert.Equal(t, userID, txn.UserID)
	assert.NotNil(t, txn.OrderID)
	assert.Equal(t, orderID, *txn.OrderID)
	assert.Equal(t, "refund", txn.TransactionType)
	assert.Equal(t, refundAmount, txn.Amount)
	assert.Equal(t, initialBalance.Balance, txn.BalanceBefore)
	assert.Equal(t, initialBalance.Balance+refundAmount, txn.BalanceAfter)
	assert.Equal(t, "completed", txn.Status)
	assert.Equal(t, "balance", txn.PaymentMethod)
	assert.Equal(t, description, txn.Description)
	assert.NotZero(t, txn.CreatedAt)
	assert.NotNil(t, txn.CompletedAt)

	// Verify balance was updated
	newBalance, err := GetUserBalance(userID)
	require.NoError(t, err)
	assert.Equal(t, initialBalance.Balance+refundAmount, newBalance.Balance)

	// Cleanup - withdraw the refunded amount
	_ = ProcessWithdrawal(userID, refundAmount, "test", "cleanup", "Test cleanup")
}

// TestProcessWithdrawal tests removing funds from user balance
func TestProcessWithdrawal(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// First deposit funds to ensure sufficient balance
	depositAmount := 500.00
	err := ProcessDeposit(userID, depositAmount, "test", "TEST"+uuid.New().String()[:8], "Setup funds")
	require.NoError(t, err)

	// Get balance after deposit
	balanceAfterDeposit, err := GetUserBalance(userID)
	require.NoError(t, err)

	withdrawalAmount := 200.00
	paymentMethod := "bank_transfer"
	referenceNumber := "BT" + uuid.New().String()[:8]
	description := "Test withdrawal"

	// Process withdrawal
	err = ProcessWithdrawal(userID, withdrawalAmount, paymentMethod, referenceNumber, description)

	require.NoError(t, err)

	// Verify balance was updated
	newBalance, err := GetUserBalance(userID)
	require.NoError(t, err)
	assert.Equal(t, balanceAfterDeposit.Balance-withdrawalAmount, newBalance.Balance)

	// Cleanup - withdraw remaining deposit
	_ = ProcessWithdrawal(userID, depositAmount-withdrawalAmount, "test", "cleanup", "Test cleanup")
}

// TestProcessWithdrawal_InsufficientFunds tests withdrawal with insufficient balance
func TestProcessWithdrawal_InsufficientFunds(t *testing.T) {
	_ = SetupDatabase()

	// Create a new user with zero balance
	newUserID := "543255dd-5325-4d3f-bcd2-ee6f8ac87e2e"
	err := CreateUserBalance(newUserID)
	require.NoError(t, err)

	// Try to withdraw more than available balance
	err = ProcessWithdrawal(newUserID, 100.00, "test", "test", "Test withdrawal")

	require.Error(t, err)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "UserBalance" WHERE user_id = $1`, newUserID)
}

// TestGetPaymentTransaction tests retrieving a specific transaction
func TestGetPaymentTransaction(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Create a transaction
	depositAmount := 300.00
	err := ProcessDeposit(userID, depositAmount, "test", "TEST"+uuid.New().String()[:8], "Test deposit")
	require.NoError(t, err)

	// Get the most recent transaction for this user
	transactions, err := GetUserTransactions(userID, 1, 0)
	require.NoError(t, err)
	require.NotEmpty(t, transactions)

	// Retrieve the transaction by ID
	retrievedTxn, err := GetPaymentTransaction(transactions[0].ID)

	require.NoError(t, err)
	assert.Equal(t, transactions[0].ID, retrievedTxn.ID)
	assert.Equal(t, userID, retrievedTxn.UserID)
	assert.Equal(t, "deposit", retrievedTxn.TransactionType)
	assert.Equal(t, depositAmount, retrievedTxn.Amount)

	// Cleanup
	_ = ProcessWithdrawal(userID, depositAmount, "test", "cleanup", "Test cleanup")
}

// TestGetPaymentTransaction_NotFound tests retrieving a non-existent transaction
func TestGetPaymentTransaction_NotFound(t *testing.T) {
	_ = SetupDatabase()

	_, err := GetPaymentTransaction("00000000-0000-0000-0000-000000000000")

	require.Error(t, err)
}

// TestGetUserTransactions tests retrieving transaction history for a user
func TestGetUserTransactions(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Create multiple transactions
	err := ProcessDeposit(userID, 100.00, "test", uuid.NewString(), "Test deposit 1")
	require.NoError(t, err)
	err = ProcessDeposit(userID, 200.00, "test", uuid.NewString(), "Test deposit 2")
	require.NoError(t, err)

	// Retrieve transactions
	transactions, err := GetUserTransactions(userID, 10, 0)

	require.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.GreaterOrEqual(t, len(transactions), 2, "Should have at least 2 transactions")

	// Verify transactions are ordered by created_at DESC
	if len(transactions) >= 2 {
		assert.True(t, transactions[0].CreatedAt.After(transactions[1].CreatedAt) ||
			transactions[0].CreatedAt.Equal(transactions[1].CreatedAt))
	}

	// Verify all transactions belong to the user
	for _, txn := range transactions {
		assert.Equal(t, userID, txn.UserID)
		assert.NotEmpty(t, txn.ID)
		assert.NotEmpty(t, txn.TransactionType)
	}

	// Cleanup
	_ = ProcessWithdrawal(userID, 300.00, "test", "cleanup", "Test cleanup")
}

// TestGetUserTransactions_Pagination tests pagination of transaction history
func TestGetUserTransactions_Pagination(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Create multiple transactions
	for i := 1; i <= 5; i++ {
		err := ProcessDeposit(userID, float64(i*50), "test", uuid.New().String()[:8], "Test deposit")
		require.NoError(t, err)
	}

	// Test pagination - first page
	page1, err := GetUserTransactions(userID, 2, 0)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(page1), 2)

	// Test pagination - second page
	page2, err := GetUserTransactions(userID, 2, 2)
	require.NoError(t, err)
	assert.NotNil(t, page2)

	// Verify pages don't overlap
	if len(page1) > 0 && len(page2) > 0 {
		assert.NotEqual(t, page1[0].ID, page2[0].ID)
	}

	// Cleanup
	_ = ProcessWithdrawal(userID, 750.00, "test", "cleanup", "Test cleanup")
}

// TestPaymentWorkflow tests a complete payment workflow
func TestPaymentWorkflow(t *testing.T) {
	_ = SetupDatabase()

	// Use existing test user JohnDoe
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Step 1: Get initial balance
	initialBalance, err := GetUserBalance(userID)
	require.NoError(t, err)

	// Step 2: Deposit funds
	err = ProcessDeposit(userID, 1000.00, "gcash", uuid.NewString(), "Initial deposit")
	require.NoError(t, err)

	// Step 3: Make a payment
	orderID := "983a8b8f-2fcf-45ce-b26f-5965ae5fef8b"
	paymentTxn, err := ProcessOrderPayment(userID, orderID, 400.00, "Order payment")
	require.NoError(t, err)
	assert.Equal(t, "order_payment", paymentTxn.TransactionType)

	// Step 4: Process a refund (e.g., order cancelled)
	refundTxn, err := ProcessRefund(userID, orderID, 400.00, "Order cancelled")
	require.NoError(t, err)
	assert.Equal(t, "refund", refundTxn.TransactionType)

	// Step 5: Withdraw funds
	err = ProcessWithdrawal(userID, 500.00, "bank", uuid.NewString(), "Withdrawal to bank")
	require.NoError(t, err)

	// Step 6: Verify final balance
	finalBalance, err := GetUserBalance(userID)
	require.NoError(t, err)
	expectedBalance := initialBalance.Balance + 1000.00 - 500.00
	assert.Equal(t, expectedBalance, finalBalance.Balance)

	// Step 7: Verify transaction history
	transactions, err := GetUserTransactions(userID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(transactions), 4, "Should have at least 4 transactions")

	// Cleanup - restore initial balance
	_ = ProcessWithdrawal(userID, 500.00, "test", "cleanup", "Test cleanup")
}
