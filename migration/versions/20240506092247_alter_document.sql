-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents ADD COLUMN IF NOT EXISTS status VARCHAR(100) NOT NULL DEFAULT 'new';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
