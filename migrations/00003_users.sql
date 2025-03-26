-- +goose Up
-- +goose StatementBegin
create table if not exists users (
	id uuid not null unique,
	email varchar(255) not null unique,
	password_hash varchar(255) not null,
	type varchar(64) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
