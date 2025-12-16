-- +goose Up
-- +goose StatementBegin
CREATE TABLE event
(
    id          uuid default gen_random_uuid() not null primary key,
    user_id     integer                        not null,
    title       text                           not null,
    description text,
    datetime    timestamp                      not null,
    duration    varchar(255),
    remind_time varchar(255),
    created_at  timestamp                      not null,
    updated_at  timestamp                      not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
