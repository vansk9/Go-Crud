CREATE TABLE cart_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  user_id BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  name TEXT,
  quantity INT,
  price DOUBLE PRECISION,
  color TEXT,
  size TEXT
);
