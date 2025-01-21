-- name: CreateWallet :one
INSERT INTO wallets (user_id)
VALUES ($1)
RETURNING user_id, balance, created_at;

-- name: GetWalletByUserId :one
SELECT user_id, balance, created_at, updated_at
FROM wallets
WHERE user_id = $1;

-- name: CreditWalletBalance :one
UPDATE wallets
SET balance = balance + $1, updated_at = NOW()
WHERE user_id = $2
RETURNING user_id;

-- name: DebitWalletBalance :one
UPDATE wallets
SET balance = balance - $1, updated_at = NOW()
WHERE user_id = $2
RETURNING user_id;