package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/moderations"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func Test_PostModerateTask(t *testing.T) {
	t.Run("Name", func(t *testing.T) {
		task := PostModerateTask{}

		require.Equal(t, name, task.Name())
	})

	t.Run("Priority", func(t *testing.T) {
		task := PostModerateTask{}

		require.NotEqual(t, 0, task.Priority())
	})

	t.Run("Validate", func(t *testing.T) {
		task := PostModerateTask{}

		require.Error(t, task.Validate())
	})

	t.Run("Args", func(t *testing.T) {
		task := PostModerateTask{
			UserID: uuid.UUID("userID"),
			PostID: 12,
			Reason: "some-reason",
		}

		require.JSONEq(t, `{
      "user-id": "userID",
      "post-id": 12,
      "reason": "some-reason"
      }`, string(task.Args()))
	})
}

func Test_PostModerateTaskRunner(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Run with an invalid json", func(t *testing.T) {
		t.Parallel()

		usersSvc := users.NewMockService(t)
		postsSvc := posts.NewMockService(t)
		moderationsSvc := moderations.NewMockService(t)
		svc := NewPostModerateTaskRunner(usersSvc, postsSvc, moderationsSvc)

		err := svc.Run(ctx, json.RawMessage(`some-invalid json`))
		require.ErrorContains(t, err, "failed to unmarshal the args")
	})

	t.Run("RunArgs success", func(t *testing.T) {
		t.Parallel()

		usersSvc := users.NewMockService(t)
		postsSvc := posts.NewMockService(t)
		moderationsSvc := moderations.NewMockService(t)
		svc := NewPostModerateTaskRunner(usersSvc, postsSvc, moderationsSvc)

		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()
		moderation := moderations.NewFakeModeration(t).Build()
		reason := gofakeit.LoremIpsumSentence(5)

		usersSvc.On("GetByID", ctx, user.ID()).Return(user, nil).Once()
		postsSvc.On("GetByID", ctx, post.ID()).Return(post, nil).Once()
		moderationsSvc.On("ModeratePost", ctx, &moderations.PostModerationCmd{
			User:   user,
			Post:   post,
			Reason: reason,
		}).Return(moderation, nil).Once()
		postsSvc.On("SetPostStatus", ctx, post, posts.Moderated).Return(nil).Once()

		err := svc.RunArgs(ctx, &PostModerateTask{
			UserID: user.ID(),
			PostID: post.ID(),
			Reason: reason,
		})
		require.NoError(t, err)
	})

	t.Run("RunArgs with a users.GetByID error", func(t *testing.T) {
		t.Parallel()

		usersSvc := users.NewMockService(t)
		postsSvc := posts.NewMockService(t)
		moderationsSvc := moderations.NewMockService(t)
		svc := NewPostModerateTaskRunner(usersSvc, postsSvc, moderationsSvc)

		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()
		reason := gofakeit.LoremIpsumSentence(5)

		usersSvc.On("GetByID", ctx, user.ID()).Return(nil, errs.Internal(errors.New("some-error"))).Once()

		err := svc.RunArgs(ctx, &PostModerateTask{
			UserID: user.ID(),
			PostID: post.ID(),
			Reason: reason,
		})
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
	})

	t.Run("RunArgs with a post.GetByID error", func(t *testing.T) {
		t.Parallel()

		usersSvc := users.NewMockService(t)
		postsSvc := posts.NewMockService(t)
		moderationsSvc := moderations.NewMockService(t)
		svc := NewPostModerateTaskRunner(usersSvc, postsSvc, moderationsSvc)

		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()
		reason := gofakeit.LoremIpsumSentence(5)

		usersSvc.On("GetByID", ctx, user.ID()).Return(user, nil).Once()
		postsSvc.On("GetByID", ctx, post.ID()).Return(nil, errs.Internal(errors.New("some-error"))).Once()

		err := svc.RunArgs(ctx, &PostModerateTask{
			UserID: user.ID(),
			PostID: post.ID(),
			Reason: reason,
		})
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
	})

	t.Run("RunArgs with a moderation.ModeratePost error", func(t *testing.T) {
		t.Parallel()

		usersSvc := users.NewMockService(t)
		postsSvc := posts.NewMockService(t)
		moderationsSvc := moderations.NewMockService(t)
		svc := NewPostModerateTaskRunner(usersSvc, postsSvc, moderationsSvc)

		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()
		reason := gofakeit.LoremIpsumSentence(5)

		usersSvc.On("GetByID", ctx, user.ID()).Return(user, nil).Once()
		postsSvc.On("GetByID", ctx, post.ID()).Return(post, nil).Once()
		moderationsSvc.On("ModeratePost", ctx, &moderations.PostModerationCmd{
			User:   user,
			Post:   post,
			Reason: reason,
		}).Return(nil, errs.Internal(errors.New("some-error"))).Once()

		err := svc.RunArgs(ctx, &PostModerateTask{
			UserID: user.ID(),
			PostID: post.ID(),
			Reason: reason,
		})
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
	})

	t.Run("RunArgs with a posts.SetPostStatus error", func(t *testing.T) {
		t.Parallel()

		usersSvc := users.NewMockService(t)
		postsSvc := posts.NewMockService(t)
		moderationsSvc := moderations.NewMockService(t)
		svc := NewPostModerateTaskRunner(usersSvc, postsSvc, moderationsSvc)

		user := users.NewFakeUser(t).Build()
		post := posts.NewFakePost(t).Build()
		moderation := moderations.NewFakeModeration(t).Build()
		reason := gofakeit.LoremIpsumSentence(5)

		usersSvc.On("GetByID", ctx, user.ID()).Return(user, nil).Once()
		postsSvc.On("GetByID", ctx, post.ID()).Return(post, nil).Once()
		moderationsSvc.On("ModeratePost", ctx, &moderations.PostModerationCmd{
			User:   user,
			Post:   post,
			Reason: reason,
		}).Return(moderation, nil).Once()
		postsSvc.On("SetPostStatus", ctx, post, posts.Moderated).Return(errs.Internal(errors.New("some-error"))).Once()

		err := svc.RunArgs(ctx, &PostModerateTask{
			UserID: user.ID(),
			PostID: post.ID(),
			Reason: reason,
		})
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
	})
}
