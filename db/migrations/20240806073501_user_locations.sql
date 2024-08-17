-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_locations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    location GEOGRAPHY(Point, 4326),
    tracked Boolean DEFAULT FALSE,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_locations;
-- +goose StatementEnd
