-- +goose Up
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(10),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE roles;