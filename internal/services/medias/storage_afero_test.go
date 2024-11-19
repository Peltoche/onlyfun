package medias

import (
	"io"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Storage_Afero(t *testing.T) {
	t.Parallel()

	t.Run("Upload and Download Delete success", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewToolboxForTest(t)
		fs := afero.NewMemMapFs()
		storage, err := newStorageAfero(fs, "/", tools)

		rawFile := []byte("Hello, World!")

		// Run
		fileID, writer, err := storage.NewFileUploader()
		require.NoError(t, err)
		require.NotEmpty(t, fileID)
		require.NotNil(t, writer)

		writer.Write(rawFile)
		writer.Close()

		// Run 2
		reader, err := storage.NewFileDownloader(fileID)
		require.NoError(t, err)
		require.NotNil(t, reader)

		// Asserts 2
		res, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, rawFile, res)
	})

	t.Run("Upload with a fs error", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewToolboxForTest(t)
		fs := afero.NewMemMapFs()
		storage, err := newStorageAfero(fs, "/", tools)
		require.NoError(t, err)

		// Create an fs error by removing the write permission
		storage.fs = afero.NewReadOnlyFs(afero.NewMemMapFs())

		fileID, reader, err := storage.NewFileUploader()
		require.Empty(t, fileID)
		require.Nil(t, reader)
		require.ErrorContains(t, err, "operation not permitted")
		require.ErrorContains(t, err, "failed to create")
	})
}
