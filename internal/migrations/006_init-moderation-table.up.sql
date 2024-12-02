CREATE TABLE IF NOT EXISTS moderations(
  "id" INTEGER PRIMARY KEY,
  "post_id" INTEGER NOT NULL,
  "reason" TEXT NOT NULL,
  "created_at" TEXT NOT NULL,
  "created_by" TEXT NOT NULL,
  FOREIGN KEY(created_by) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  FOREIGN KEY(post_id) REFERENCES posts(id) ON UPDATE RESTRICT ON DELETE RESTRICT
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_id ON moderations(id);
CREATE INDEX IF NOT EXISTS idx_moderations_created_by ON moderations(created_by);
