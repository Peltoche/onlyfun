package moderations

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func Test_Moderations_Service(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("ModeratePost success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, permsSvc)

		now := time.Now()
		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()
		reason := gofakeit.LoremIpsumSentence(10)
		moderation := Moderation{
			id:        3232,
			postID:    post.ID(),
			reason:    reason,
			createdAt: now,
			createdBy: user.ID(),
		}
		moderationWithoutID := moderation
		moderationWithoutID.id = 0

		permsSvc.On("IsAuthorized", user, perms.Moderation).Return(true).Once()
		tools.ClockMock.On("Now").Return(now).Once()
		storage.On("Save", ctx, &moderationWithoutID).Return(nil).Once()

		res, err := svc.ModeratePost(ctx, &PostModerationCmd{
			User:   user,
			Post:   post,
			Reason: reason,
		})
		require.NoError(t, err)
		require.Equal(t, &moderationWithoutID, res)
	})

	t.Run("ModeratePost with an invalid authorization error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, permsSvc)

		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()

		permsSvc.On("IsAuthorized", user, perms.Moderation).Return(false).Once()

		res, err := svc.ModeratePost(ctx, &PostModerationCmd{
			User:   user,
			Post:   post,
			Reason: gofakeit.LoremIpsumSentence(10),
		})
		require.Nil(t, res)
		require.ErrorIs(t, err, errs.ErrUnauthorized)
	})

	t.Run("ModeratePost with a validation error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, permsSvc)

		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()

		permsSvc.On("IsAuthorized", user, perms.Moderation).Return(true).Once()

		res, err := svc.ModeratePost(ctx, &PostModerationCmd{
			User:   user,
			Post:   post,
			Reason: "f",
		})
		require.Nil(t, res)
		require.ErrorIs(t, err, errs.ErrValidation)
	})

	t.Run("ModeratePost with a Save error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		permsSvc := perms.NewMockService(t)
		svc := newService(tools, storage, permsSvc)

		now := time.Now()
		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()
		reason := gofakeit.LoremIpsumSentence(10)
		moderation := Moderation{
			id:        3232,
			postID:    post.ID(),
			reason:    reason,
			createdAt: now,
			createdBy: user.ID(),
		}
		moderationWithoutID := moderation
		moderationWithoutID.id = 0

		permsSvc.On("IsAuthorized", user, perms.Moderation).Return(true).Once()
		tools.ClockMock.On("Now").Return(now).Once()
		storage.On("Save", ctx, &moderationWithoutID).Return(errors.New("some-error")).Once()

		res, err := svc.ModeratePost(ctx, &PostModerationCmd{
			User:   user,
			Post:   post,
			Reason: reason,
		})
		require.Nil(t, res)
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
	})
}
