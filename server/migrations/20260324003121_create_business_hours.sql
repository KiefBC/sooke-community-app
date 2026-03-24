-- +goose Up
-- +goose StatementBegin
CREATE TABLE business_hours (
  id BIGSERIAL PRIMARY KEY,
  business_id BIGINT NOT NULL REFERENCES businesses (id) ON DELETE CASCADE,
  day_of_week SMALLINT NOT NULL CHECK (
    day_of_week >= 0
    AND day_of_week <= 6
  ),
  open_time TIME NOT NULL,
  close_time TIME NOT NULL,
  is_closed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  UNIQUE (business_id, day_of_week) -- will ensure that there is only one entry per business per day of the week
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS business_hours;

-- +goose StatementEnd
