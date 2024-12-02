package posts

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/roles"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/require"
)

func TestPostSqlStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("GetListedPosts with nothing", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		// Run
		res, err := store.GetListedPosts(ctx, 0, 10)

		// Asserts
		require.NoError(t, err)
		require.Empty(t, res)
	})

	t.Run("GetLatestPostWithStatus with nothing", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		res, err := store.GetLatestPostWithStatus(ctx, Listed)

		require.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("Save success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role := roles.NewFakeRole(t).BuildAndStore(ctx, db)
		avatar := medias.NewFakeFileMeta(t).BuildAndStore(ctx, db)
		user := users.NewFakeUser(t).WithRole(role).WithAvatar(avatar).BuildAndStore(ctx, db)
		post := NewFakePost(t).
			CreatedBy(user).
			Build()

		oldPostID := post.ID()

		// Run
		err := store.Save(ctx, post)

		// Asserts
		require.NoError(t, err)
		require.NotEqual(t, oldPostID, post.ID())
	})

	t.Run("GetLatestPostWithStatus success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role := roles.NewFakeRole(t).BuildAndStore(ctx, db)
		avatar := medias.NewFakeFileMeta(t).BuildAndStore(ctx, db)
		user := users.NewFakeUser(t).WithRole(role).WithAvatar(avatar).BuildAndStore(ctx, db)
		post := NewFakePost(t).CreatedBy(user).WithStatus(Listed).BuildAndStore(ctx, db)
		_ = NewFakePost(t).CreatedBy(user).WithStatus(Uploaded).BuildAndStore(ctx, db)

		res, err := store.GetLatestPostWithStatus(ctx, Listed)

		require.NoError(t, err)
		require.Equal(t, post, res)
	})

	t.Run("GetListedPosts success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		nbUsers := 25
		role := roles.NewFakeRole(t).BuildAndStore(ctx, db)
		avatar := medias.NewFakeFileMeta(t).BuildAndStore(ctx, db)
		user := users.NewFakeUser(t).WithRole(role).WithAvatar(avatar).BuildAndStore(ctx, db)
		posts := make(map[uint64]Post, nbUsers)

		for i := 0; i < nbUsers; i++ {
			res := NewFakePost(t).CreatedBy(user).WithStatus(Listed).BuildAndStore(ctx, db)
			t.Logf("%02d -> %s\n", res.id, res.title)
			posts[res.id] = *res
		}

		// Test 1
		res, err := store.GetListedPosts(ctx, 24, 5)
		require.NoError(t, err)
		require.EqualValues(t, []Post{
			posts[24],
			posts[23],
			posts[22],
			posts[21],
			posts[20],
		}, res)

		// Test 2
		res2, err2 := store.GetListedPosts(ctx, 4, 10)
		require.NoError(t, err2)
		require.EqualValues(t, []Post{
			posts[4],
			posts[3],
			posts[2],
			posts[1],
		}, res2)
	})

	t.Run("GetLatestPostWithStatus and GetListedPosts success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		nbUsers := 25
		role := roles.NewFakeRole(t).BuildAndStore(ctx, db)
		avatar := medias.NewFakeFileMeta(t).BuildAndStore(ctx, db)
		user := users.NewFakeUser(t).WithRole(role).WithAvatar(avatar).BuildAndStore(ctx, db)
		posts := make(map[uint64]Post, nbUsers)

		for i := 0; i < nbUsers; i++ {
			res := NewFakePost(t).CreatedBy(user).WithStatus(Listed).BuildAndStore(ctx, db)
			t.Logf("%02d -> %s\n", res.id, res.title)
			posts[res.id] = *res
		}

		// Test 1
		latest, err := store.GetLatestPostWithStatus(ctx, Listed)
		require.NoError(t, err)

		res, err := store.GetListedPosts(ctx, latest.id, 5)
		require.NoError(t, err)
		require.EqualValues(t, []Post{
			posts[25],
			posts[24],
			posts[23],
			posts[22],
			posts[21],
		}, res)
	})

	t.Run("CountPostsWithStatus success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		nbListed := 5
		nbUploaded := 5

		role := roles.NewFakeRole(t).BuildAndStore(ctx, db)
		avatar := medias.NewFakeFileMeta(t).BuildAndStore(ctx, db)
		user := users.NewFakeUser(t).WithRole(role).WithAvatar(avatar).BuildAndStore(ctx, db)

		for i := 0; i < nbListed; i++ {
			_ = NewFakePost(t).CreatedBy(user).WithStatus(Listed).BuildAndStore(ctx, db)
		}

		for i := 0; i < nbUploaded; i++ {
			_ = NewFakePost(t).CreatedBy(user).WithStatus(Uploaded).BuildAndStore(ctx, db)
		}

		// Test 1
		resListed, err := store.CountPostsWithStatus(ctx, Listed)
		require.NoError(t, err)
		require.Equal(t, nbListed, resListed)

		// Test 2
		resUploaded, err := store.CountPostsWithStatus(ctx, Listed)
		require.NoError(t, err)
		require.Equal(t, nbUploaded, resUploaded)
	})
}
