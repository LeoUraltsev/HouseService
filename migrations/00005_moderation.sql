-- +goose Up
-- +goose StatementBegin
create table if not exists moderation (
    moderator_id uuid not null,
    flat_id BIGINT unique not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists moderation;
-- +goose StatementEnd
