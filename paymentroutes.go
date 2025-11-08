package main

import (
	"bukidlink/db"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Balance Management Endpoints

// getUserBalanceHandler retrieves the balance for a specific user
func getUserBalanceHandler(c *gin.Context) {
	userID := c.Param("user_id")

	balance, err := db.GetUserBalance(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User Balance not found"})
		}
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, balance)
}

// createUserBalanceHandler initializes a balance record for a new user
func createUserBalanceHandler(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	err := db.CreateUserBalance(req.UserID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Balance created successfully", "user_id": req.UserID})
}

// Deposit & Withdrawal Endpoints

// processDepositHandler adds funds to user's balance
func processDepositHandler(c *gin.Context) {
	var req struct {
		UserID          string  `json:"user_id" binding:"required"`
		Amount          float64 `json:"amount" binding:"required,gt=0"`
		PaymentMethod   string  `json:"payment_method" binding:"required"`
		ReferenceNumber string  `json:"reference_number" binding:"required"`
		Description     string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	err := db.ProcessDeposit(req.UserID, req.Amount, req.PaymentMethod, req.ReferenceNumber, req.Description)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit processed successfully"})
}

// processWithdrawalHandler removes funds from user balance
func processWithdrawalHandler(c *gin.Context) {
	var req struct {
		UserID          string  `json:"user_id" binding:"required"`
		Amount          float64 `json:"amount" binding:"required,gt=0"`
		PaymentMethod   string  `json:"payment_method" binding:"required"`
		ReferenceNumber string  `json:"reference_number" binding:"required"`
		Description     string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	err := db.ProcessWithdrawal(req.UserID, req.Amount, req.PaymentMethod, req.ReferenceNumber, req.Description)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal processed successfully"})
}

// Order Payments & Refunds Endpoints

// processOrderPaymentHandler deducts funds from user balance for an order
func processOrderPaymentHandler(c *gin.Context) {
	var req struct {
		UserID      string  `json:"user_id" binding:"required"`
		OrderID     string  `json:"order_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	txn, err := db.ProcessOrderPayment(req.UserID, req.OrderID, req.Amount, req.Description)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, txn)
}

// processRefundHandler returns funds to user balance
func processRefundHandler(c *gin.Context) {
	var req struct {
		UserID      string  `json:"user_id" binding:"required"`
		OrderID     string  `json:"order_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	txn, err := db.ProcessRefund(req.UserID, req.OrderID, req.Amount, req.Description)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, txn)
}

// Transaction History Endpoints

// getPaymentTransactionHandler retrieves a specific transaction by ID
func getPaymentTransactionHandler(c *gin.Context) {
	transactionID := c.Param("transaction_id")

	txn, err := db.GetPaymentTransaction(transactionID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, txn)
}

// getUserTransactionsHandler retrieves paginated transaction history for a user
func getUserTransactionsHandler(c *gin.Context) {
	userID := c.Param("user_id")

	// Parse pagination parameters
	limit := 20 // default
	offset := 0 // default

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	transactions, err := db.GetUserTransactions(userID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No transactions found"})
			return
		}
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, transactions)
}
