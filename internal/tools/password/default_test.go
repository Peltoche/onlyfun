package password

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("Encrypt/Compare success", func(t *testing.T) {
		password := &Argon2IDPassword{}

		hashed, err := password.Encrypt(ctx, secret.NewText("some-password"))
		require.NoError(t, err)
		require.NotEqual(t, "some-password", hashed)

		ok, err := password.Compare(ctx, hashed, secret.NewText("some-password"))
		assert.True(t, ok)
		require.NoError(t, err)
	})

	t.Run("Decrypte with a no base64 string", func(t *testing.T) {
		password := &Argon2IDPassword{}

		ok, err := password.Compare(ctx, secret.NewText("not a hex string#"), secret.NewText("some-password"))
		assert.False(t, ok)
		require.EqualError(t, err, "failed to decode the hash: the encoded hash is not in the correct format")
	})
}
