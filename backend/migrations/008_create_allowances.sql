-- +goose Up
CREATE TABLE IF NOT EXISTS allowances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    period_start DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_allowances_user_id ON allowances(user_id);

-- +goose Down
DROP TABLE IF EXISTS allowances;
