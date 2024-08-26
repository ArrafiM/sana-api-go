-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS merchandise_images (
    id BIGSERIAL PRIMARY KEY,
    url VARCHAR(255) NOT NULL,
    merchandise_id BIGINT NOT NULL,
    FOREIGN KEY (merchandise_id) REFERENCES merchandises(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS merchandise_images;
-- +goose StatementEnd
