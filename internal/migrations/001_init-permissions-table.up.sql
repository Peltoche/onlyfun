CREATE TABLE IF NOT EXISTS permissions (
  "role" TEXT NOT NULL,
  "permissions" TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_role_name ON permissions(role);
