CREATE TABLE group_status (
  group_id UUID,
  delta integer,
  total integer,
  created_at timestamp with time zone
);

CREATE INDEX group_status_group_id_idx on group_status(group_id,created_at);