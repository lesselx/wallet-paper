CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    wallet_id bigint NOT NULL REFERENCES wallets(id),
    amount numeric(15, 2) NOT NULL,
    transaction_type text NOT NULL CHECK (transaction_type IN ('topup', 'withdraw')),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);