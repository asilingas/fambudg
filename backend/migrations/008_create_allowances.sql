-- +goose Up
CREATE TABLE IF NOT EXISTS allowances (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    period_start DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_allowances_uuid ON allowances(uuid);
CREATE INDEX idx_allowances_user_id ON allowances(user_id);

-- +goose Down
DROP TABLE IF EXISTS allowances;
