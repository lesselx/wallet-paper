CREATE TABLE IF NOT EXISTS wallets (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id),
    balance numeric(15, 2) NOT NULL DEFAULT 0.00,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
