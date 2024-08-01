-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD picture VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN picture;
-- +goose StatementEnd
