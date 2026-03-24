-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
  id BIGSERIAL PRIMARY KEY,
  event_type_id BIGINT NOT NULL REFERENCES event_types (id),
  submitted_by BIGINT NOT NULL REFERENCES users (id),
  business_id BIGINT REFERENCES businesses (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  description TEXT,
  latitude DOUBLE PRECISION,
  longitude DOUBLE PRECISION,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ,
  status TEXT NOT NULL DEFAULT 'draft' CHECK (
    status IN (
      'draft',
      'pending_review',
      'approved',
      'cancelled',
      'rejected'
    )
  ),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  -- the constraint ensures that either business_id is provided (and lat/long are null) or lat/long are provided (and business_id is null)
  CONSTRAINT chk_event_location CHECK (
    (
      business_id IS NOT NULL
      AND latitude IS NULL
      AND longitude IS NULL
    )
    OR (
      business_id IS NULL
      AND latitude IS NOT NULL
      AND longitude IS NOT NULL
    )
  )
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;

-- +goose StatementEnd
