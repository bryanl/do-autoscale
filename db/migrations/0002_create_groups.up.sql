CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE groups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  base_name text,
  base_size int,
  metric_type text,
  template_id int references templates(id)
);