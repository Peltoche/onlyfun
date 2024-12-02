package moderations

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/require"
)

func Test_Moderations_SqlStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("GetAll with nothing", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		// Run
		res, err := store.GetAll(ctx, &sqlstorage.PaginateCmd{Limit: 10})

		// Asserts
		require.NoError(t, err)
		require.Empty(t, res)
	})

	t.Run("Save success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role, _ := perms.NewFakePermissions(t).BuildAndStore(ctx, db)
		avatar := medias.NewFakeFileMeta(t).BuildAndStore(ctx, db)
		user := users.NewFakeUser(t).WithRole(role).WithAvatar(avatar).BuildAndStore(ctx, db)
		post := posts.NewFakePost(t).CreatedBy(user).BuildAndStore(ctx, db)
		moderation := NewFakeModeration(t).CreatedBy(user).WithPost(post).Build()

		// Run
		err := store.Save(ctx, moderation)

		// Asserts
		require.NoError(t, err)
		require.NotEqual(t, uint64(0), post.ID())
	})
}
