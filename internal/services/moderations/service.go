package moderations

import (
	"context"
	"fmt"

	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

type storage interface {
	Save(ctx context.Context, m *Moderation) error
	GetAll(ctx context.Context, cmd *sqlstorage.PaginateCmd) ([]Moderation, error)
}

type service struct {
	clock    clock.Clock
	uuid     uuid.Service
	permsSvc perms.Service
	storage  storage
}

func newService(tools tools.Tools, storage storage, permsSvc perms.Service) *service {
	svc := &service{
		clock:    tools.Clock(),
		uuid:     tools.UUID(),
		storage:  storage,
		permsSvc: permsSvc,
	}

	return svc
}

func (s *service) ModeratePost(ctx context.Context, cmd *PostModerationCmd) (*Moderation, error) {
	if !s.permsSvc.IsAuthorized(cmd.User, perms.Moderation) {
		return nil, errs.Unauthorized(fmt.Errorf("user %q doesn't have the authorization %q", cmd.User.ID(), perms.Moderation))
	}

	err := cmd.Validate()
	if err != nil {
		return nil, errs.Validation(err)
	}

	moderation := Moderation{
		// id: set by the db
		postID:    cmd.Post.ID(),
		reason:    cmd.Reason,
		createdAt: s.clock.Now(),
		createdBy: cmd.User.ID(),
	}

	err = s.storage.Save(ctx, &moderation)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to Save in db: %w", err))
	}

	return &moderation, nil
}
