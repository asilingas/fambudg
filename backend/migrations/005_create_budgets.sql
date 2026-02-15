-- +goose Up
CREATE TABLE budgets (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    amount BIGINT NOT NULL,
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    year INT NOT NULL CHECK (year >= 2000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (category_id, month, year)
);

CREATE INDEX idx_budgets_uuid ON budgets(uuid);
CREATE INDEX idx_budgets_category_id ON budgets(category_id);
CREATE INDEX idx_budgets_month_year ON budgets(month, year);

-- +goose Down
DROP TABLE IF EXISTS budgets;
