-- +goose Up
-- +goose StatementBegin
ALTER TABLE merchants ADD
picture VARCHAR(255)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE merchants
DROP COLUMN picture;
-- +goose StatementEnd
