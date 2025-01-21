-- name: CreateWallet :one
INSERT INTO wallets (user_id)
VALUES ($1)
RETURNING user_id, balance, created_at;

