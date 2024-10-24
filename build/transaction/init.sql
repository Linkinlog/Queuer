CREATE TABLE IF NOT EXISTS transactions (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL,
  result JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX transactions_id_idx ON transactions (id);
