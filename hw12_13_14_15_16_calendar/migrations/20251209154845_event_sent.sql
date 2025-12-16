-- +goose Up
-- +goose StatementBegin
ALTER TABLE event ADD COLUMN IF NOT EXISTS remind_sent_time timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event DROP COLUMN IF EXISTS remind_sent_time;
-- +goose StatementEnd
