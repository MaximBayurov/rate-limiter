-- +goose Up
-- +goose StatementBegin
CREATE TABLE ip_list (
    ip INET NOT NULL UNIQUE,
    type VARCHAR(10) CHECK (type IN ('white', 'black'))
);

-- Индексы для оптимизации запросов
CREATE INDEX index_ip_list_ip ON ip_list(ip);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ip_list
-- +goose StatementEnd
