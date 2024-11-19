CREATE TABLE IF NOT EXISTS medias (
  "id" TEXT NOT NULL,
  "size" INTEGER NOT NULL,
  "mimetype" TEXT NOT NULL,
  "type" TEXT NOT NULL,
  "checksum" TEXT NOT NULL,
  "uploaded_at" TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_medias_id ON medias(id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_medias_checksum ON medias(checksum);
