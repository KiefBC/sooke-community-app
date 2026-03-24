-- +goose Up
-- +goose StatementBegin
CREATE TABLE menus (
  id BIGSERIAL PRIMARY KEY,
  business_id BIGINT NOT NULL REFERENCES businesses (id) ON DELETE CASCADE,
  name TEXT NOT NULL, -- lunch menu, breakfast menu, etc
  description TEXT, -- maybe dave's breakfast menu has a description that says "the best breakfast in town"
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS menus;

-- +goose StatementEnd
