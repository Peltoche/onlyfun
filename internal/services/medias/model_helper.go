package medias

import (
	"context"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type FakeFileMetaBuilder struct {
	t        testing.TB
	fileMeta *FileMeta
}

func NewFakeFileMeta(t testing.TB) *FakeFileMetaBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()

	return &FakeFileMetaBuilder{
		t: t,
		fileMeta: &FileMeta{
			id:         uuidProvider.New(),
			mimetype:   "image/jpeg",
			checksum:   "AVQ5FFAO1rvc9OD6bYPUHxLWS9AZZ2/u4Rl2fYEwID8",
			mediaType:  Post,
			size:       1024,
			uploadedAt: gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now()),
		},
	}
}

func (f *FakeFileMetaBuilder) Build() *FileMeta {
	return f.fileMeta
}
func (f *FakeFileMetaBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *FileMeta {
	f.t.Helper()

	storage := newSqlStorage(db)

	fileMeta := f.Build()

	err := storage.Save(ctx, fileMeta)
	require.NoError(f.t, err)

	return fileMeta
}
