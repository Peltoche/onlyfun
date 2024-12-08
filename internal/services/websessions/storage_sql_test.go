package websessions

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionSqlStorage(t *testing.T) {
	db := sqlstorage.NewTestStorage(t)
	storage := newSQLStorage(db)
	ctx := context.Background()

	role, _ := perms.NewFakePermissions(t).BuildAndStore(ctx, db)
	avatar := medias.NewFakeFileMeta(t).BuildAndStore(ctx, db)
	user := users.NewFakeUser(t).WithRole(role).WithAvatar(avatar).BuildAndStore(ctx, db)
	sessionToken := "some-token"
	session := NewFakeSession(t).
		CreatedBy(user).
		WithToken(sessionToken).
		Build()

	t.Run("Create success", func(t *testing.T) {
		// Run
		err := storage.Save(context.Background(), session)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByToken success", func(t *testing.T) {
		// Run
		res, err := storage.GetByToken(context.Background(), secret.NewText(sessionToken))

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, session, res)
	})

	t.Run("GeAllForUser success", func(t *testing.T) {
		// Run
		res, err := storage.GetAllForUser(context.Background(), user.ID(), nil)

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, []Session{*session}, res)
	})

	t.Run("GetByToken not found", func(t *testing.T) {
		// Run
		res, err := storage.GetByToken(context.Background(), secret.NewText("some-invalid-token"))

		// Asserts
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("RemoveByToken ", func(t *testing.T) {
		// Run
		err := storage.RemoveByToken(context.Background(), secret.NewText(sessionToken))

		// Asserts
		require.NoError(t, err)
		res, err := storage.GetByToken(context.Background(), secret.NewText(sessionToken))
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})
}
