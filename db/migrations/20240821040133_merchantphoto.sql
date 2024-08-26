-- +goose Up
-- +goose StatementBegin
ALTER TABLE merchants 
ADD picture VARCHAR(255), 
ADD color VARCHAR(7);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE merchants
DROP COLUMN picture, 
DROP COLUMN color;
-- +goose StatementEnd
