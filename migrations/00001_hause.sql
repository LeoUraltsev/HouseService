-- +goose Up
-- +goose StatementBegin
create table if not exists house (
	id BIGINT unique primary key generated always as identity,
	address TEXT not null,
	year SMALLINT not null,
	developer TEXT,
	created_at TIMESTAMP not null default now(),
	last_flat_add_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exist hause;
-- +goose StatementEnd
