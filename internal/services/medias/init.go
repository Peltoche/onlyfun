package medias

import (
	"context"
	"fmt"
	"io"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/spf13/afero"
)

type Service interface {
	Upload(ctx context.Context, mediaType MediaType, r io.Reader) (*FileMeta, error)
	Download(ctx context.Context, id uuid.UUID) (io.ReadSeekCloser, error)
	GetMetadataByChecksum(ctx context.Context, checksum string) (*FileMeta, error)
	GetMetadata(ctx context.Context, fileID uuid.UUID) (*FileMeta, error)
	Delete(ctx context.Context, fileID uuid.UUID) error
}

func Init(
	dirPath string,
	fs afero.Fs,
	tools tools.Tools,
	db sqlstorage.Querier,
) (Service, error) {
	fileStorage, err := newStorageAfero(fs, dirPath, tools)
	if err != nil {
		return nil, fmt.Errorf("failed to setup the afero storage: %w", err)
	}

	mediaStorage := newSqlStorage(db)

	return newService(fileStorage, mediaStorage, tools), nil
}
