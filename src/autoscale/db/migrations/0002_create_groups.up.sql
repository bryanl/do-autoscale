CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE groups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name text not null,
  base_name text,
  template_name text references templates(name),
  metric_type text,
  metric jsonb,
  policy_type text,
  policy jsonb
);

CREATE UNIQUE INDEX groups_name_idx on groups(name);