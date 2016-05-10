ALTER TABLE templates
  ADD name text;

UPDATE templates
set name = id
where name is null;

ALTER TABLE templates
ALTER COLUMN name set not null;

ALTER TABLE templates
ADD CONSTRAINT templates_name_key UNIQUE(name);