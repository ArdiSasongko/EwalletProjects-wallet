-- name: InsertWalletTransactionsCredit :one
INSERT INTO wallet_transactions (wallet_id, amount, wallet_transactions_type, reference)
VALUES ($1, $2, 'credit', $3)
RETURNING wallet_id, reference, amount, created_at;

-- name: InsertWalletTransactionsDebit :one
INSERT INTO wallet_transactions (wallet_id, amount, wallet_transactions_type, reference)
VALUES ($1, $2, 'debit', $3)
RETURNING wallet_id, reference, amount, created_at;

-- name: GetHistoryTransactions :many
SELECT wallet_id, reference, amount, wallet_transactions_type, created_at
FROM wallet_transactions
WHERE wallet_id = $1 AND wallet_transactions_type = $2
ORDER BY created_at DESC
LIMIT $3
OFFSET $4;