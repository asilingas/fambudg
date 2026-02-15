-- +goose Up
CREATE TABLE saving_goals (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    target_amount BIGINT NOT NULL,
    current_amount BIGINT NOT NULL DEFAULT 0,
    target_date DATE,
    priority INT NOT NULL DEFAULT 1,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_saving_goals_uuid ON saving_goals(uuid);
CREATE INDEX idx_saving_goals_status ON saving_goals(status);
CREATE INDEX idx_saving_goals_priority ON saving_goals(priority);

-- +goose Down
DROP TABLE IF EXISTS saving_goals;
