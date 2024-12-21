package posts

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/stretchr/testify/require"
)

func Test_Posts_Service(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Create success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		mediaContent := strings.NewReader("some-content")

		fileMeta := medias.NewFakeFileMeta(t).Build()
		user := users.NewFakeUser(t).Build()
		post := NewFakePost(t).CreatedBy(user).WithMedia(fileMeta).Build()

		postWithoutID := post
		postWithoutID.id = 0

		mediasSvc.On("Upload", ctx, medias.Post, mediaContent).Return(fileMeta, nil).Once()
		tools.ClockMock.On("Now").Return(post.CreatedAt).Once()
		storage.On("Save", ctx, postWithoutID).Return(nil).Once()

		res, err := svc.Create(ctx, &CreateCmd{
			Title:     post.title,
			Media:     mediaContent,
			CreatedBy: user,
		})
		require.NoError(t, err)
		require.Equal(t, post, res)
	})

	t.Run("Create with a validation error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		user := users.NewFakeUser(t).Build()
		mediaContent := strings.NewReader("some-content")

		res, err := svc.Create(ctx, &CreateCmd{
			Title:     "fo", // Too short
			Media:     mediaContent,
			CreatedBy: user,
		})
		require.ErrorIs(t, err, errs.ErrValidation)
		require.Nil(t, res)
	})

	t.Run("Create with a media.Upload error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		mediaContent := strings.NewReader("some-content")
		user := users.NewFakeUser(t).Build()

		mediasSvc.On("Upload", ctx, medias.Post, mediaContent).Return(nil, fmt.Errorf("some-error")).Once()

		res, err := svc.Create(ctx, &CreateCmd{
			Title:     "Some title",
			Media:     mediaContent,
			CreatedBy: user,
		})
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		require.Nil(t, res)
	})

	t.Run("Create with a Save error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		mediaContent := strings.NewReader("some-content")

		fileMeta := medias.NewFakeFileMeta(t).Build()
		user := users.NewFakeUser(t).Build()
		post := NewFakePost(t).CreatedBy(user).WithMedia(fileMeta).Build()

		postWithoutID := post
		postWithoutID.id = 0

		mediasSvc.On("Upload", ctx, medias.Post, mediaContent).Return(fileMeta, nil).Once()
		tools.ClockMock.On("Now").Return(post.CreatedAt).Once()
		storage.On("Save", ctx, postWithoutID).Return(fmt.Errorf("some-error")).Once()

		res, err := svc.Create(ctx, &CreateCmd{
			Title:     post.title,
			Media:     mediaContent,
			CreatedBy: user,
		})
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		require.Nil(t, res)
	})

	t.Run("GetByID success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		post := NewFakePost(t).Build()

		storage.On("GetByID", ctx, post.ID()).Return(post, nil).Once()

		res, err := svc.GetByID(ctx, post.ID())
		require.NoError(t, err)
		require.Equal(t, post, res)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("GetByID", ctx, uint(32)).Return(nil, errNotFound).Once()

		res, err := svc.GetByID(ctx, 32)
		require.ErrorIs(t, err, errs.ErrNotFound)
		require.Nil(t, res)
	})

	t.Run("GetByID with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("GetByID", ctx, uint(32)).Return(nil, fmt.Errorf("some-error")).Once()

		res, err := svc.GetByID(ctx, 32)
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		require.Nil(t, res)
	})

	t.Run("GetNextPostToModerate success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		post := NewFakePost(t).Build()

		storage.On("GetOldestPostWithStatus", ctx, Uploaded).Return(post, nil).Once()

		res, err := svc.GetNextPostToModerate(ctx)
		require.NoError(t, err)
		require.Equal(t, post, res)
	})

	t.Run("GetNextPostToModerate not fund", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("GetOldestPostWithStatus", ctx, Uploaded).Return(nil, errNotFound).Once()

		res, err := svc.GetNextPostToModerate(ctx)
		require.ErrorIs(t, err, errs.ErrNotFound)
		require.Nil(t, res)
	})

	t.Run("GetNextPostToModerate with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("GetOldestPostWithStatus", ctx, Uploaded).Return(nil, fmt.Errorf("some-error")).Once()

		res, err := svc.GetNextPostToModerate(ctx)
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		require.Nil(t, res)
	})

	t.Run("GetUserStats success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		user := users.NewFakeUser(t).Build()

		storage.On("CountUserPostsByStatus", ctx, user.ID(), Uploaded).Return(12, nil).Once()
		storage.On("CountUserPostsByStatus", ctx, user.ID(), Listed).Return(23, nil).Once()
		storage.On("CountUserPostsByStatus", ctx, user.ID(), Moderated).Return(5, nil).Once()

		res, err := svc.GetUserStats(ctx, user)
		require.NoError(t, err)
		require.Equal(t, map[Status]int{
			Uploaded:  12,
			Listed:    23,
			Moderated: 5,
		}, res)
	})

	t.Run("GetUserStats with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		user := users.NewFakeUser(t).Build()

		storage.On("CountUserPostsByStatus", ctx, user.ID(), Uploaded).Return(12, nil).Once()
		storage.On("CountUserPostsByStatus", ctx, user.ID(), Listed).Return(0, fmt.Errorf("some-error")).Once()

		res, err := svc.GetUserStats(ctx, user)
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		require.Nil(t, res)
	})

	t.Run("GetLatestPost success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		post := NewFakePost(t).Build()

		storage.On("GetLatestPostWithStatus", ctx, Listed).Return(post, nil).Once()

		res, err := svc.GetLatestPost(ctx)
		require.NoError(t, err)
		require.Equal(t, post, res)
	})

	t.Run("GetLatestPost not found", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("GetLatestPostWithStatus", ctx, Listed).Return(nil, errNotFound).Once()

		res, err := svc.GetLatestPost(ctx)
		require.ErrorIs(t, err, errs.ErrNotFound)
		require.Nil(t, res)
	})

	t.Run("GetLatestPost with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("GetLatestPostWithStatus", ctx, Listed).Return(nil, fmt.Errorf("some-error")).Once()

		res, err := svc.GetLatestPost(ctx)
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		require.Nil(t, res)
	})

	t.Run("CountPostsWaitingModeration success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("CountPostsWithStatus", ctx, Uploaded).Return(32, nil).Once()

		res, err := svc.CountPostsWaitingModeration(ctx)
		require.NoError(t, err)
		require.Equal(t, 32, res)
	})

	t.Run("CountPostsWaitingModeration with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("CountPostsWithStatus", ctx, Uploaded).Return(0, fmt.Errorf("some-error")).Once()

		res, err := svc.CountPostsWaitingModeration(ctx)
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		require.Equal(t, 0, res)
	})

	t.Run("GetPosts success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		posts := make([]Post, 3)
		for i := range 3 {
			posts[i] = *NewFakePost(t).Build()
		}

		storage.On("GetListedPosts", ctx, uint(200), uint(3)).Return(posts, nil).Once()

		res, err := svc.GetPosts(ctx, 200, 3)
		require.NoError(t, err)
		require.Equal(t, posts, res)
	})

	t.Run("GetPosts with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, mediasSvc, permsSvc)

		storage.On("GetListedPosts", ctx, uint(200), uint(3)).Return(nil, fmt.Errorf("some-error")).Once()

		res, err := svc.GetPosts(ctx, 200, 3)
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		require.Nil(t, res)
	})
}
