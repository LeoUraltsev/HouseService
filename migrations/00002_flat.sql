-- +goose Up
-- +goose StatementBegin
create table if not exists flat (
	id SERIAL primary key not null,
	house_id BIGINT not null,
	price BIGINT not null,
	rooms SMALLINT not null,
	status VARCHAR(255) not null,
	foreign key (house_id) references house (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exist flat;
-- +goose StatementEnd
