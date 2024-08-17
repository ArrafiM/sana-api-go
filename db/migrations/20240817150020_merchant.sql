-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS merchants (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(50),
    description TEXT,
    status VARCHAR(10) DEFAULT 'active',
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE merchants;
-- +goose StatementEnd
