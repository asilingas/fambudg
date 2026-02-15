-- +goose Up
CREATE TABLE bill_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    amount BIGINT NOT NULL,
    due_day INT NOT NULL CHECK (due_day BETWEEN 1 AND 31),
    frequency VARCHAR(20) NOT NULL DEFAULT 'monthly' CHECK (frequency IN ('monthly', 'quarterly', 'yearly')),
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    next_due_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bill_reminders_next_due_date ON bill_reminders(next_due_date);
CREATE INDEX idx_bill_reminders_is_active ON bill_reminders(is_active);

-- +goose Down
DROP TABLE IF EXISTS bill_reminders;
