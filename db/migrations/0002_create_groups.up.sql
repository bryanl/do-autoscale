CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE groups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name text not null,
  base_name text,
  base_size int,
  metric_type text,
  template_name text references templates(name)
);

CREATE UNIQUE INDEX groups_name_idx on groups(name);