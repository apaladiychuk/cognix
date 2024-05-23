-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents ADD COLUMN IF NOT EXISTS status VARCHAR(100) NOT NULL DEFAULT 'new';
ALTER TABLE documents DROP COLUMN IF EXISTS boost;
ALTER TABLE documents DROP COLUMN IF EXISTS hidden;
ALTER TABLE documents DROP COLUMN IF EXISTS semantic_id;
ALTER TABLE documents DROP COLUMN IF EXISTS from_ingestion_api;
ALTER TABLE connectors DROP COLUMN IF EXISTS input_type;
    -- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
