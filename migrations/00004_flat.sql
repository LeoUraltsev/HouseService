-- +goose Up
-- +goose StatementBegin
alter table if exists flat add if not exists number BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists flat drop number if exists;
-- +goose StatementEnd
