-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks ( -- создание таблицы задач
    id SERIAL PRIMARY KEY, -- id задачи
    title VARCHAR(255) NOT NULL, -- название задачи
    completed BOOLEAN NOT NULL DEFAULT FALSE -- выполнена ли задача
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks; -- удаление таблицы задач если она существует
-- +goose StatementEnd
