BEGIN;

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP(6);
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TABLE IF NOT EXISTS registry (
  id bigserial PRIMARY KEY,
  network text NOT NULL,
  address text NOT NULL,
  name text NOT NULL,
  created_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
  deleted_at timestamp(6) NULL DEFAULT NULL
);

CREATE TRIGGER update_updated_at_registry
BEFORE UPDATE ON flix
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE TABLE IF NOT EXISTS flix (
  id bigserial PRIMARY KEY,
  flix_id text NOT NULL,
  json_body JSONB NOT NULL,
  created_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
  deleted_at timestamp(6) NULL DEFAULT NULL,
);

CREATE TRIGGER update_updated_at_flix
BEFORE UPDATE ON flix
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE TABLE IF NOT EXISTS flix_hash (
  id bigserial PRIMARY KEY,
  flix_id text NOT NULL,
  registry_id bigserial NOT NULL,
  cadence_body_hash text NOT NULL,
  CONSTRAINT flix_hash_uk_flix_id UNIQUE (flix_id)
  CONSTRAINT flix_hash_uk_registry_id UNIQUE (registry_id)
);

CREATE TABLE IF NOT EXISTS registry_flix (
  id bigserial PRIMARY KEY,
  flix_id text NOT NULL,
  registry_id bigserial NOT NULL,
  CONSTRAINT registry_flix_uk_flix_id UNIQUE (flix_id)
  CONSTRAINT registry_flix_uk_registry_id UNIQUE (registry_id)
);

CREATE TABLE IF NOT EXISTS alias (
  id bigserial PRIMARY KEY,
  name text NOT NULL,
  flix_id bigserial NOT NULL,
  created_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
  deleted_at timestamp(6) NULL DEFAULT NULL,
  CONSTRAINT alias_uk_flix_id UNIQUE (flix_id)
);

CREATE TRIGGER update_updated_at_alias
BEFORE UPDATE ON alias
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE TABLE IF NOT EXISTS auditor (
  id bigserial PRIMARY KEY,
  name text NOT NULL,
  address text NOT NULL,
  x_url text NULL DEFAULT NULL,
  website_url text NULL DEFAULT NULL,
  created_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
  deleted_at timestamp(6) NULL DEFAULT NULL,
  CONSTRAINT alias_uk_template_id UNIQUE (template_id)
);

CREATE TRIGGER update_updated_at_alias
BEFORE UPDATE ON alias
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

COMMIT;
