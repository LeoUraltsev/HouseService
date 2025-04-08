-- +goose Up
-- +goose StatementBegin
create table if not exists flat (
	id BIGINT GENERATED ALWAYS AS IDENTITY,
	house_id BIGINT not null,
	price BIGINT not null,
	rooms SMALLINT not null,
	status VARCHAR(255) not null,
	foreign key (house_id) references house (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists flat;
-- +goose StatementEnd
