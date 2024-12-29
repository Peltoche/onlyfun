package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Peltoche/onlyfun/internal/services/moderations"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const name = "post-moderate"

type PostModerateTask struct {
	UserID uuid.UUID `json:"user-id"`
	PostID uint      `json:"post-id"`
	Reason string    `json:"reason"`
}

func (r *PostModerateTask) Name() string  { return name }
func (r *PostModerateTask) Priority() int { return 1 }

func (r *PostModerateTask) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.UserID, v.Required, is.UUIDv4),
		v.Field(&r.PostID, v.Required, is.UUIDv4),
		v.Field(&r.Reason, v.Required, v.Length(3, 1000)),
	)
}

func (r *PostModerateTask) Args() json.RawMessage {
	res, _ := json.Marshal(r)

	return res
}

type PostModerateTaskRunner struct {
	usersSvc       users.Service
	postsSvc       posts.Service
	moderationsSvc moderations.Service
}

func NewPostModerateTaskRunner(
	usersSvc users.Service,
	postsSvc posts.Service,
	moderationsSvc moderations.Service,
) *PostModerateTaskRunner {
	return &PostModerateTaskRunner{
		usersSvc:       usersSvc,
		postsSvc:       postsSvc,
		moderationsSvc: moderationsSvc,
	}
}

func (r *PostModerateTaskRunner) Name() string { return name }

func (r *PostModerateTaskRunner) Run(ctx context.Context, rawArgs json.RawMessage) error {
	var args PostModerateTask

	err := json.Unmarshal(rawArgs, &args)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the args: %w", err)
	}

	return r.RunArgs(ctx, &args)
}

func (r *PostModerateTaskRunner) RunArgs(ctx context.Context, args *PostModerateTask) error {
	user, err := r.usersSvc.GetByID(ctx, args.UserID)
	if err != nil {
		return fmt.Errorf("failed to get the user %q: %w", args.UserID, err)
	}

	post, err := r.postsSvc.GetByID(ctx, args.PostID)
	if err != nil {
		return fmt.Errorf("failed to get the post %q: %w", args.PostID, err)
	}

	_, err = r.moderationsSvc.ModeratePost(ctx, &moderations.PostModerationCmd{
		User:   user,
		Post:   post,
		Reason: args.Reason,
	})
	if err != nil {
		return fmt.Errorf("failed to moderate post %w", err)
	}

	err = r.postsSvc.SetPostStatus(ctx, post, posts.Moderated)
	if err != nil {
		return fmt.Errorf("failed to SetPostStatus: %w", err)
	}

	return nil
}
