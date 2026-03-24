-- +goose Up
-- +goose StatementBegin
CREATE TABLE menu_items (
  id BIGSERIAL PRIMARY KEY,
  menu_id BIGINT NOT NULL REFERENCES menus (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  price NUMERIC(10, 2) NOT NULL, -- numeric with 10 digits total and 2 decimal places, so max price is 99999999.99
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX idx_menu_items_menu_id ON menu_items (menu_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS menu_items;

-- +goose StatementEnd
