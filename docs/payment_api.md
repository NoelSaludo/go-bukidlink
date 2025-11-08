# BukidLink Payment & Wallet API Documentation

## Overview
The Payment & Wallet system provides in-app balance management, deposits, withdrawals, order payments, and refunds. All monetary amounts are in Philippine Peso (PHP). The system uses database transactions with row-level locking (`FOR UPDATE`) to prevent race conditions during concurrent balance updates.

## Base URL
All endpoints are prefixed with `/balance` or `/payment`

---

## Balance Management

### Get User Balance
Retrieves the current balance for a specific user.

**Endpoint:** `GET /balance/:user_id`

**Path Parameters:**
- `user_id` (string, required): UUID of the user

**Response (200 OK):**
```json
{
  "id": "bal-uuid",
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
  "balance": 5000.00,
  "currency": "PHP",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-15T14:30:00Z"
}
```

**Error Responses:**
- `404 Not Found`: User balance not found
- `500 Internal Server Error`: Database error

**Example:**
```bash
curl -X GET http://localhost:8080/balance/d30869ec-fb97-46d8-85a3-82608c01f803
```

---

### Create User Balance
Initializes a balance record for a new user with 0.00 PHP starting balance.

**Endpoint:** `POST /balance`

**Request Body:**
```json
{
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803"
}
```

**Response (201 Created):**
```json
{
  "message": "Balance created successfully",
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803"
}
```

**Error Responses:**
- `400 Bad Request`: Missing or invalid user_id
- `409 Conflict`: Balance already exists for this user
- `500 Internal Server Error`: Database error

**Example:**
```bash
curl -X POST http://localhost:8080/balance \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803"
  }'
```

---

## Deposit & Withdrawal

### Process Deposit
Adds funds to a user's balance. Creates a completed transaction record and updates the balance atomically.

**Endpoint:** `POST /payment/deposit`

**Request Body:**
```json
{
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
  "amount": 1000.00,
  "payment_method": "gcash",
  "reference_number": "GC2025110800001",
  "description": "Top-up via GCash"
}
```

**Fields:**
- `user_id` (string, required): UUID of the user
- `amount` (float, required): Amount to deposit (must be > 0)
- `payment_method` (string, required): Payment method used (e.g., "gcash", "maya", "bank_transfer")
- `reference_number` (string, required): External transaction reference (unique)
- `description` (string, optional): Transaction description

**Response (200 OK):**
```json
{
  "message": "Deposit processed successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or amount ≤ 0
- `404 Not Found`: User balance not found
- `409 Conflict`: Duplicate reference number
- `500 Internal Server Error`: Database error

**Database Function:** `ProcessDeposit(userID, amount, paymentMethod, referenceNumber, description) error`

**Implementation Details:**
1. Begins database transaction
2. Locks user balance row with `SELECT ... FOR UPDATE`
3. Calculates new balance: `balanceAfter = balanceBefore + amount`
4. Inserts transaction record with type `'deposit'` and status `'completed'`
5. Updates user balance
6. Commits transaction

**Example:**
```bash
curl -X POST http://localhost:8080/payment/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "amount": 1000.00,
    "payment_method": "gcash",
    "reference_number": "GC2025110800001",
    "description": "Top-up via GCash"
  }'
```

---

### Process Withdrawal
Removes funds from a user's balance. Validates sufficient balance before processing.

**Endpoint:** `POST /payment/withdrawal`

**Request Body:**
```json
{
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
  "amount": 500.00,
  "payment_method": "bank_transfer",
  "reference_number": "BNK2025110800001",
  "description": "Withdraw to BDO account"
}
```

**Fields:**
- `user_id` (string, required): UUID of the user
- `amount` (float, required): Amount to withdraw (must be > 0)
- `payment_method` (string, required): Withdrawal method (e.g., "bank_transfer", "gcash", "maya")
- `reference_number` (string, required): External transaction reference (unique)
- `description` (string, optional): Transaction description

**Response (200 OK):**
```json
{
  "message": "Withdrawal processed successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or amount ≤ 0
- `404 Not Found`: User balance not found or insufficient funds
- `409 Conflict`: Duplicate reference number
- `500 Internal Server Error`: Database error

**Database Function:** `ProcessWithdrawal(userID, amount, paymentMethod, referenceNumber, description) error`

**Implementation Details:**
1. Begins database transaction
2. Locks user balance row with `SELECT ... FOR UPDATE`
3. Validates sufficient balance: `balance.Balance >= amount`
4. Calculates new balance: `balanceAfter = balanceBefore - amount`
5. Inserts transaction record with type `'withdrawal'` and status `'completed'`
6. Updates user balance
7. Commits transaction
8. Returns `sql.ErrNoRows` if insufficient funds

**Example:**
```bash
curl -X POST http://localhost:8080/payment/withdrawal \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "amount": 500.00,
    "payment_method": "bank_transfer",
    "reference_number": "BNK2025110800001",
    "description": "Withdraw to BDO account"
  }'
```

---

## Order Payments & Refunds

### Process Order Payment
Deducts funds from user balance for an order. Validates order existence and sufficient balance.

**Endpoint:** `POST /payment/order`

**Request Body:**
```json
{
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
  "order_id": "11111111-1111-1111-1111-111111111111",
  "amount": 450.50,
  "description": "Payment for order - Fresh vegetables"
}
```

**Fields:**
- `user_id` (string, required): UUID of the user
- `order_id` (string, required): UUID of the order being paid
- `amount` (float, required): Payment amount (must be > 0)
- `description` (string, optional): Payment description

**Response (200 OK):**
```json
{
  "id": "txn-uuid",
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
  "order_id": "11111111-1111-1111-1111-111111111111",
  "transaction_type": "order_payment",
  "amount": 450.50,
  "balance_before": 5000.00,
  "balance_after": 4549.50,
  "status": "completed",
  "payment_method": "balance",
  "reference_number": null,
  "description": "Payment for order - Fresh vegetables",
  "created_at": "2024-01-15T14:30:00Z",
  "updated_at": "2024-01-15T14:30:00Z",
  "completed_at": "2024-01-15T14:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or amount ≤ 0
- `404 Not Found`: Order not found, user balance not found, or insufficient funds
- `500 Internal Server Error`: Database error

**Database Function:** `ProcessOrderPayment(userID, orderID, amount, description) (PaymentTransaction, error)`

**Implementation Details:**
1. Begins database transaction
2. Validates order exists: `SELECT id FROM "Order" WHERE id = $1`
3. Locks user balance row with `SELECT ... FOR UPDATE`
4. Validates sufficient balance: `balance.Balance >= amount`
5. Calculates new balance: `balanceAfter = balanceBefore - amount`
6. Inserts transaction record with type `'order_payment'`, status `'completed'`, and payment_method `'balance'`
7. Uses `RETURNING` clause to retrieve full transaction details
8. Updates user balance
9. Commits transaction
10. Returns complete PaymentTransaction object

**Example:**
```bash
curl -X POST http://localhost:8080/payment/order \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "order_id": "11111111-1111-1111-1111-111111111111",
    "amount": 450.50,
    "description": "Payment for order - Fresh vegetables"
  }'
```

---

### Process Refund
Returns funds to user balance for a refunded order.

**Endpoint:** `POST /payment/refund`

**Request Body:**
```json
{
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
  "order_id": "11111111-1111-1111-1111-111111111111",
  "amount": 450.50,
  "description": "Refund for cancelled order"
}
```

**Fields:**
- `user_id` (string, required): UUID of the user
- `order_id` (string, required): UUID of the order being refunded
- `amount` (float, required): Refund amount (must be > 0)
- `description` (string, optional): Refund description

**Response (200 OK):**
```json
{
  "id": "txn-uuid",
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
  "order_id": "11111111-1111-1111-1111-111111111111",
  "transaction_type": "refund",
  "amount": 450.50,
  "balance_before": 4549.50,
  "balance_after": 5000.00,
  "status": "completed",
  "payment_method": "balance",
  "reference_number": null,
  "description": "Refund for cancelled order",
  "created_at": "2024-01-15T15:00:00Z",
  "updated_at": "2024-01-15T15:00:00Z",
  "completed_at": "2024-01-15T15:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or amount ≤ 0
- `404 Not Found`: User balance not found
- `500 Internal Server Error`: Database error

**Database Function:** `ProcessRefund(userID, orderID, amount, description) (PaymentTransaction, error)`

**Implementation Details:**
1. Begins database transaction
2. Locks user balance row with `SELECT ... FOR UPDATE`
3. Calculates new balance: `balanceAfter = balanceBefore + amount`
4. Inserts transaction record with type `'refund'`, status `'completed'`, and payment_method `'balance'`
5. Uses `RETURNING` clause to retrieve full transaction details
6. Updates user balance
7. Commits transaction
8. Returns complete PaymentTransaction object

**Example:**
```bash
curl -X POST http://localhost:8080/payment/refund \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "order_id": "11111111-1111-1111-1111-111111111111",
    "amount": 450.50,
    "description": "Refund for cancelled order"
  }'
```

---

## Transaction History

### Get Payment Transaction
Retrieves a specific transaction by its ID.

**Endpoint:** `GET /payment/transaction/:transaction_id`

**Path Parameters:**
- `transaction_id` (string, required): UUID of the transaction

**Response (200 OK):**
```json
{
  "id": "20000001-0000-0000-0000-000000000001",
  "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
  "order_id": null,
  "transaction_type": "deposit",
  "amount": 5000.00,
  "balance_before": 0.00,
  "balance_after": 5000.00,
  "status": "completed",
  "payment_method": "gcash",
  "reference_number": "GC2025110800001",
  "description": "Initial deposit via GCash",
  "created_at": "2025-10-30T08:30:00+08:00",
  "updated_at": "2025-10-30T08:30:00+08:00",
  "completed_at": "2025-10-30T08:30:00+08:00"
}
```

**Error Responses:**
- `404 Not Found`: Transaction not found
- `500 Internal Server Error`: Database error

**Database Function:** `GetPaymentTransaction(transactionID) (PaymentTransaction, error)`

**Example:**
```bash
curl -X GET http://localhost:8080/payment/transaction/20000001-0000-0000-0000-000000000001
```

---

### Get User Transactions
Retrieves paginated transaction history for a user, ordered by creation date (newest first).

**Endpoint:** `GET /payment/transactions/:user_id`

**Path Parameters:**
- `user_id` (string, required): UUID of the user

**Query Parameters:**
- `limit` (int, optional): Number of transactions to return (default: 20, must be > 0)
- `offset` (int, optional): Number of transactions to skip for pagination (default: 0, must be >= 0)

**Response (200 OK):**
```json
[
  {
    "id": "20000010-0000-0000-0000-000000000010",
    "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "order_id": null,
    "transaction_type": "order_payment",
    "amount": 1000.00,
    "balance_before": 6000.00,
    "balance_after": 5000.00,
    "status": "completed",
    "payment_method": "balance",
    "reference_number": null,
    "description": "Payment for order - Fresh fruits",
    "created_at": "2025-11-07T14:30:00+08:00",
    "updated_at": "2025-11-07T14:30:00+08:00",
    "completed_at": "2025-11-07T14:30:00+08:00"
  },
  {
    "id": "20000009-0000-0000-0000-000000000009",
    "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "order_id": null,
    "transaction_type": "deposit",
    "amount": 1000.00,
    "balance_before": 5000.00,
    "balance_after": 6000.00,
    "status": "completed",
    "payment_method": "maya",
    "reference_number": "MAYA2025110800002",
    "description": "Top-up via Maya",
    "created_at": "2025-11-07T08:00:00+08:00",
    "updated_at": "2025-11-07T08:00:00+08:00",
    "completed_at": "2025-11-07T08:00:00+08:00"
  },
  {
    "id": "20000001-0000-0000-0000-000000000001",
    "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "order_id": null,
    "transaction_type": "deposit",
    "amount": 5000.00,
    "balance_before": 0.00,
    "balance_after": 5000.00,
    "status": "completed",
    "payment_method": "gcash",
    "reference_number": "GC2025110800001",
    "description": "Initial deposit via GCash",
    "created_at": "2025-10-30T08:30:00+08:00",
    "updated_at": "2025-10-30T08:30:00+08:00",
    "completed_at": "2025-10-30T08:30:00+08:00"
  }
]
```

**Error Responses:**
- `500 Internal Server Error`: Database error

**Database Function:** `GetUserTransactions(userID, limit, offset) ([]PaymentTransaction, error)`

**Query:** Orders transactions by `created_at DESC` (newest first)

**Notes:**
- Returns empty array `[]` if user has no transactions
- Invalid limit/offset values (non-numeric, negative, zero) fall back to defaults (20/0)
- Handler validates and sanitizes pagination parameters before calling database function

**Example:**
```bash
# Get first 10 transactions
curl -X GET "http://localhost:8080/payment/transactions/d30869ec-fb97-46d8-85a3-82608c01f803?limit=10&offset=0"

# Get next 10 transactions (pagination)
curl -X GET "http://localhost:8080/payment/transactions/d30869ec-fb97-46d8-85a3-82608c01f803?limit=10&offset=10"

# Get all transactions with default pagination (20 items)
curl -X GET http://localhost:8080/payment/transactions/d30869ec-fb97-46d8-85a3-82608c01f803
```

---

## Transaction Types

| Type | Description | Affects Balance | Returns Transaction Object |
|------|-------------|-----------------|----------------------------|
| `deposit` | User adds funds via external payment method | Increases (+) | No (only error) |
| `withdrawal` | User withdraws funds to external account | Decreases (-) | No (only error) |
| `order_payment` | User pays for an order using balance | Decreases (-) | Yes |
| `refund` | Order cancelled/refunded, funds returned | Increases (+) | Yes |

---

## Transaction Status Values

| Status | Description | Current Usage |
|--------|-------------|---------------|
| `completed` | Transaction successfully processed | All current transactions |
| `pending` | Transaction initiated but not yet processed | Future use (async payments) |
| `processing` | Transaction being processed | Future use (async payments) |
| `failed` | Transaction failed | Future use (async payments) |
| `cancelled` | Transaction cancelled | Future use |

**Note:** All transactions in the current implementation are immediately set to `completed` status.

---

## Data Models

### UserBalance
```go
type UserBalance struct {
    ID        string    `json:"id"`         // UUID (Primary Key)
    UserID    string    `json:"user_id"`    // Foreign key to User (Unique)
    Balance   float64   `json:"balance"`    // Current balance (>= 0.00)
    Currency  string    `json:"currency"`   // Always "PHP"
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Constraints:**
- `balance >= 0` (check constraint)
- `user_id` is UNIQUE
- Cascade delete on User removal

**Database Table:** `"UserBalance"`

---

### PaymentTransaction
```go
type PaymentTransaction struct {
    ID              string     `json:"id"`                         // UUID (Primary Key)
    UserID          string     `json:"user_id"`                    // Foreign key to User
    OrderID         *string    `json:"order_id,omitempty"`         // Foreign key to Order (nullable)
    TransactionType string     `json:"transaction_type"`           // Enum: deposit, withdrawal, order_payment, refund
    Amount          float64    `json:"amount"`                     // Transaction amount (> 0)
    BalanceBefore   float64    `json:"balance_before"`             // Balance before transaction
    BalanceAfter    float64    `json:"balance_after"`              // Balance after transaction
    Status          string     `json:"status"`                     // Enum: pending, processing, completed, failed, cancelled
    PaymentMethod   string     `json:"payment_method,omitempty"`   // e.g., "gcash", "maya", "balance"
    ReferenceNumber *string    `json:"reference_number,omitempty"` // External payment reference (nullable, unique)
    Description     string     `json:"description,omitempty"`      // Transaction description
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
    CompletedAt     *time.Time `json:"completed_at,omitempty"`     // Timestamp when completed (nullable)
}
```

**Constraints:**
- `amount > 0` (check constraint)
- `reference_number` is UNIQUE when not NULL
- `order_id` SET NULL on Order deletion
- Cascade delete on User removal

**Database Table:** `"PaymentTransaction"`

**Indexes:**
- `idx_payment_transaction_user_id` on `user_id`
- `idx_payment_transaction_order_id` on `order_id`
- `idx_payment_transaction_status` on `status`
- `idx_payment_transaction_created_at` on `created_at DESC`
- `idx_payment_transaction_reference` on `reference_number`

---

## Database Implementation Details

### Row-Level Locking Pattern
All balance-modifying operations use pessimistic locking to prevent race conditions:

```sql
SELECT id, user_id, balance 
FROM "UserBalance" 
WHERE user_id = $1 
FOR UPDATE
```

This ensures that:
1. Only one transaction can modify a user's balance at a time
2. Other transactions wait until the lock is released
3. Balance calculations are always based on the latest committed value

### Transaction Flow Example (Deposit)

```go
func ProcessDeposit(userID string, amount float64, paymentMethod, referenceNumber, description string) error {
    tx, err := db.Begin()  // 1. Start transaction
    if err != nil {
        return err
    }

    // 2. Lock balance row
    var balance UserBalance
    err = tx.QueryRow(`SELECT id, user_id, balance FROM "UserBalance" WHERE user_id = $1 FOR UPDATE`,
        userID).Scan(&balance.ID, &balance.UserID, &balance.Balance)
    if err != nil {
        return rollbackAndReturn(tx, err)
    }

    // 3. Calculate new balance
    balanceBefore := balance.Balance
    balanceAfter := balanceBefore + amount
    transactionID := uuid.New().String()
    now := time.Now()

    // 4. Insert transaction record
    _, err = tx.Exec(`INSERT INTO "PaymentTransaction" 
        (id, user_id, transaction_type, amount, balance_before, balance_after, status, 
         payment_method, reference_number, description, created_at, updated_at, completed_at)
        VALUES ($1, $2, 'deposit', $3, $4, $5, 'completed', $6, $7, $8, $9, $9, $9)`,
        transactionID, userID, amount, balanceBefore, balanceAfter,
        paymentMethod, referenceNumber, description, now)
    if err != nil {
        return rollbackAndReturn(tx, err)
    }

    // 5. Update balance
    _, err = tx.Exec(`UPDATE "UserBalance" SET balance = $1, updated_at = $2 WHERE user_id = $3`,
        balanceAfter, now, userID)
    if err != nil {
        return rollbackAndReturn(tx, err)
    }

    // 6. Commit transaction
    return tx.Commit()
}
```

### Error Handling Strategy

**Insufficient Balance:**
```go
if balance.Balance < amount {
    return rollbackAndReturn(tx, sql.ErrNoRows)
}
```
- Returns `sql.ErrNoRows` to indicate resource not available
- API layer maps to `404 Not Found`

**Order Validation:**
```go
var existingOrderID string
err = tx.QueryRow(`SELECT id FROM "Order" WHERE id = $1`, orderID).Scan(&existingOrderID)
if err != nil {
    if err == sql.ErrNoRows {
        return txn, rollbackAndReturn(tx, fmt.Errorf("order not found: %s", orderID))
    }
    return txn, rollbackAndReturn(tx, err)
}
```
- Validates order exists before creating payment transaction
- Prevents foreign key violations with clear error messages

---

## Important Notes

1. **Concurrency Safety**: All balance-modifying operations use `SELECT ... FOR UPDATE` row-level locking to prevent race conditions when multiple transactions occur simultaneously.

2. **Transaction Atomicity**: Each operation (deposit, withdrawal, order payment, refund) runs in a database transaction that either fully succeeds or fully rolls back, ensuring balance consistency.

3. **Nullable Fields**: 
   - `order_id`: NULL for deposits/withdrawals, set for order_payment/refund
   - `reference_number`: NULL for order payments/refunds, required for deposits/withdrawals
   - `completed_at`: Can be NULL for pending transactions (currently unused)

4. **Error Handling**: 
   - Insufficient balance returns `sql.ErrNoRows` → maps to `404 Not Found`
   - Order not found returns custom error → maps to `404 Not Found`
   - Database errors return generic error → maps to `500 Internal Server Error`

5. **Foreign Key Constraints**: 
   - Order payments validate order existence before creating transaction
   - User deletion cascades to UserBalance and PaymentTransaction
   - Order deletion sets PaymentTransaction.order_id to NULL

6. **Balance Constraints**:
   - Balance must be >= 0.00 (database check constraint)
   - Transaction amount must be > 0.00 (database check constraint)
   - One balance record per user (unique constraint on user_id)

7. **Reference Number Uniqueness**:
   - External payment references must be unique (prevents duplicate deposits/withdrawals)
   - Order payments don't use reference numbers (use internal balance)

8. **Return Value Patterns**:
   - Deposits/Withdrawals return only `error` (simple operations)
   - Order Payments/Refunds return `(PaymentTransaction, error)` (need transaction details)
   - Balance queries return full `UserBalance` object
   - Transaction queries return full `PaymentTransaction` object(s)

---

## Testing

Sample test data is available in `sql/payment_schema.sql` with three test users:
- **JohnDoe** (`d30869ec-fb97-46d8-85a3-82608c01f803`): Balance 5000.00 PHP
- **DanielGaliego** (`c6554794-849f-4338-87c5-6db2e2f76514`): Balance 3500.50 PHP  
- **User 3** (`6a24dd2b-d441-4b39-ab85-8fa2bd61065e`): Balance 8250.75 PHP

Comprehensive database tests are in `db/db_payment_test.go` covering:
- Balance retrieval and creation
- Deposit/withdrawal processing
- Order payment with insufficient funds validation
- Refund processing
- Transaction history retrieval with pagination
- Complete payment workflows

API integration tests can be added to `main_test.go` following existing patterns.