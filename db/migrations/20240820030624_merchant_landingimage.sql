-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS merchant_landingimages (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    FOREIGN KEY (merchant_id) REFERENCES merchants(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS merchant_landingimages;
-- +goose StatementEnd
