-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS image (
    id SERIAL PRIMARY KEY,
    status VARCHAR(15) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS image;
-- +goose StatementEnd
