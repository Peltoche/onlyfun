CREATE TABLE IF NOT EXISTS posts (
  "id" INTEGER PRIMARY KEY,
  "status" TEXT NOT NULL,
  "title" TEXT NOT NULL,
  "file_id" TEXT NOT NULL,
  "created_at" TEXT NOT NULL,
  "created_by" TEXT NOT NULL,
  FOREIGN KEY(created_by) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_status_id ON posts(status, id);
CREATE INDEX IF NOT EXISTS idx_posts_created_by_status ON posts(created_by, status);
