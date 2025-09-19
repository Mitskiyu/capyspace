-- +goose Up
-- +goose StatementBegin
CREATE TABLE widgets (
    id UUID PRIMARY KEY,
    space_id UUID NOT NULL,
    type VARCHAR(32) NOT NULL,
    x_pos INTEGER NOT NULL,
    y_pos INTEGER NOT NULL,
    minimized BOOLEAN NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    CONSTRAINT fk_widgets_spaces FOREIGN KEY(space_id) REFERENCES spaces(id) ON DELETE CASCADE
);

CREATE INDEX idx_widgets_space_id ON widgets(space_id);

CREATE TRIGGER widgets_update_modified_at
    BEFORE UPDATE ON widgets
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS widgets_update_modified_at ON widgets;
DROP TABLE widgets;
-- +goose StatementEnd
