package medias

import (
	"testing"
)

func TestMediaStorageAfero(t *testing.T) {
	t.Parallel()

	// TODO: Write the missing tests

	t.Skip("Test not written yet")

	// t.Run("Upload with a copy error", func(t *testing.T) {
	// 	t.Parallel()
	//
	// 	tools := tools.NewToolboxForTest(t)
	// 	fs := afero.NewMemMapFs()
	// 	svc := newService(fs, tools)
	//
	// 	// Create a file
	// 	fileMeta, err := svc.Upload(ctx, iotest.ErrReader(fmt.Errorf("some-error")))
	// 	require.ErrorContains(t, err, "upload error")
	// 	require.ErrorContains(t, err, "some-error")
	// 	assert.Nil(t, fileMeta)
	// })
}
