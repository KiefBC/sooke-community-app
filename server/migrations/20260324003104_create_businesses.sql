-- +goose Up
-- +goose StatementBegin
CREATE TABLE businesses (
  id BIGSERIAL PRIMARY KEY,
  owner_id BIGINT REFERENCES users (id) ON DELETE SET NULL,
  category_id BIGINT REFERENCES business_categories (id),
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  description TEXT,
  phone TEXT,
  email TEXT,
  website TEXT,
  address TEXT NOT NULL,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX idx_businesses_owner_id ON businesses (owner_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS businesses;

-- +goose StatementEnd
