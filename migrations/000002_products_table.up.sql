CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  name TEXT,
  description TEXT,
  quantity INT,
  price DOUBLE PRECISION,
  color TEXT,
  size TEXT
);
