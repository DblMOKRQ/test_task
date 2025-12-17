CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    CONSTRAINT balance_must_be_non_negative CHECK (balance >= 0)
);
