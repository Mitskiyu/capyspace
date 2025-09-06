-- +goose Up
-- +goose StatementBegin
CREATE TABLE spaces (
    id UUID PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    is_private BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    CONSTRAINT fk_spaces_users FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER spaces_update_modified_at
    BEFORE UPDATE ON spaces
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS spaces_update_modified_at ON spaces;
DROP TABLE IF EXISTS spaces;
-- +goose StatementEnd
