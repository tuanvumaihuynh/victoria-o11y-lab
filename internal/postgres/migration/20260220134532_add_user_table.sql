-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
	id UUID PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	password_hash BYTEA NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
