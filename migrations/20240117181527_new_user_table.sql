-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS peoples
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT NOT NULL,
    country TEXT NOT NULL,
    gender TEXT NOT NULL,
    age INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS peoples;
-- +goose StatementEnd
