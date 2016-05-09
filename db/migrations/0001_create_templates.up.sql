CREATE TABLE templates (
  id serial primary key,
  region text,
  size text,
  image text,
  ssh_keys text,
  user_data text
);