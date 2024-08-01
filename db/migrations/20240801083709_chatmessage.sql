-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat_massages (
    id BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT,
    to_user_id BIGINT,
    FOREIGN KEY (from_user_id) REFERENCES users(id),
    FOREIGN KEY (to_user_id) REFERENCES users(id),
    message TEXT,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chat_massages;
-- +goose StatementEnd
