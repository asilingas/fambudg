-- +goose Up
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    parent_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(10) NOT NULL,
    icon VARCHAR(50),
    sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_categories_uuid ON categories(uuid);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_type ON categories(type);
CREATE INDEX idx_categories_sort_order ON categories(sort_order);

-- +goose Down
DROP TABLE IF EXISTS categories;
