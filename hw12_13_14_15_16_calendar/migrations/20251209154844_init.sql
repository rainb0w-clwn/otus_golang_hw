-- +goose Up
-- +goose StatementBegin
CREATE TABLE event
(
    id          uuid               default gen_random_uuid() not null primary key,
    user_id     integer   not null,
    title       text      not null,
    description text,
    datetime    timestamp not null,
    duration    varchar(255),
    remind_time timestamp,
    created_at  timestamp not null default now(),
    updated_at  timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
