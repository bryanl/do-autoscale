CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name text not null,
  region text,
  size text,
  image text,
  ssh_keys jsonb,
  user_data text
);

CREATE UNIQUE INDEX templates_name_idx on templates(name);