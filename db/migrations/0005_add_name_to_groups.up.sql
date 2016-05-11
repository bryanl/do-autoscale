ALTER TABLE groups
  ADD name text;

UPDATE groups
set name = id
where name is null;

ALTER TABLE groups
ALTER COLUMN name set not null;

ALTER TABLE groups
ADD CONSTRAINT groups_name_key UNIQUE(name);