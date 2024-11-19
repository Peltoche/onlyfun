CREATE TABLE IF NOT EXISTS roles (
  "name" TEXT NOT NULL,
  "permissions" TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_role_name ON roles(name);
