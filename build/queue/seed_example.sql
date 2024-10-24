-- Adder
INSERT INTO
  resources (name, version, schema)
VALUES
  ('adder', '1', '{"field": "array"}');

INSERT INTO
  events (resource_id, data)
VALUES
  (1, '{"addends": [1, 2, 3]}');

-- Squarer
INSERT INTO
  resources (name, version, schema)
VALUES
  ('squarer', '1', '{"field": "object"}');

INSERT INTO
	events (resource_id, data)
VALUES
	(2, '{"factor": 2, "base": 4}');

-- Longrunner
INSERT INTO
  resources (name, version, schema)
VALUES
  ('longrunner', '1', '{"field": "number"}');

INSERT INTO
  events (resource_id, data)
VALUES
  (3, '{"time_to_run": 500}');
