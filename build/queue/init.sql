CREATE TABLE IF NOT EXISTS events (
  id SERIAL PRIMARY KEY,
  resource_id INTEGER NOT NULL,
  data JSONB NOT NULL,
  processed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX events_resource_id_idx ON events (resource_id);

CREATE UNIQUE INDEX events_id_idx ON events (id);

CREATE TABLE IF NOT EXISTS resources (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  version VARCHAR(255) NOT NULL,
  schema JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX resources_name_version_idx ON resources (name, version);

ALTER TABLE events ADD FOREIGN KEY (resource_id) REFERENCES resources (id);
