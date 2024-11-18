-- +goose Up
-- +goose StatementBegin
ALTER TABLE merchandises 
ADD tag varchar[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE merchandises
DROP COLUMN tag;
-- +goose StatementEnd
