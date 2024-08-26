-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS merchandises (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50),
    description TEXT,
    price BIGINT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    merchant_id BIGINT NOT NULL,
    FOREIGN KEY (merchant_id) REFERENCES merchants(id) ON DELETE CASCADE,
    picture VARCHAR(255),
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS merchandises;
-- +goose StatementEnd
