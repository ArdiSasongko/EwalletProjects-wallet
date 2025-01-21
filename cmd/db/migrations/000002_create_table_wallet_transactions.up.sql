CREATE TYPE wallet_transactions_type AS ENUM('credit', 'debit');

CREATE TABLE IF NOT EXISTS wallet_transactions(
    id SERIAL PRIMARY KEY,
    wallet_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    wallet_transactions_type wallet_transactions_type NOT NULL,
    reference VARCHAR(100) NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_wallet_transactions_wallet_id FOREIGN KEY (wallet_id) REFERENCES wallets(user_id)
);