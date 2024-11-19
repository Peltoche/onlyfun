CREATE TABLE IF NOT EXISTS users (
  "id" TEXT NOT NULL,
  "username" TEXT NOT NULL,
  "password" TEXT NOT NULL,
  "role" TEXT NOT NULL,
  "status" TEXT NOT NULL,
  "password_changed_at" TEXT NOT NULL,
  "avatar" TEXT NOT NULL,
  "created_at" TEXT NOT NULL,
  "created_by" TEXT NOT NULL,
  FOREIGN KEY(role) REFERENCES roles(name) ON UPDATE RESTRICT ON DELETE RESTRICT
  FOREIGN KEY(avatar) REFERENCES medias(id) ON UPDATE RESTRICT ON DELETE RESTRICT
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_id ON users(id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
