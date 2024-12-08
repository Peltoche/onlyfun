package perms

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

var ErrInvalidRoleName = fmt.Errorf("invalid role name")

type storage interface {
	Save(ctx context.Context, roles *Role, perms []Permission) error
	GetAll(ctx context.Context) (map[Role][]Permission, error)
	GetPermissions(ctx context.Context, roles *Role) ([]Permission, error)
}

type service struct {
	storage storage
	clock   clock.Clock
	uuid    uuid.Service
	logger  *slog.Logger

	permsByRole map[Role][]Permission
	lock        *sync.RWMutex
}

func newService(tools tools.Tools, storage storage) *service {
	return &service{
		storage: storage,
		clock:   tools.Clock(),
		uuid:    tools.UUID(),
		logger:  tools.Logger(),

		permsByRole: map[Role][]Permission{},
		lock:        new(sync.RWMutex),
	}
}

func (s *service) bootstrap(ctx context.Context) error {
	roles, err := s.storage.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to GetAll: %w", err)
	}

	if len(roles) == 0 {
		err = s.createDefaultRoles(ctx)
		if err != nil {
			return fmt.Errorf("failed to create the default roles")
		}

		roles, err = s.storage.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("failed to GetAll after the creation: %w", err)
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for role, perms := range roles {
		s.permsByRole[role] = perms
	}

	return nil
}

func (s *service) IsAuthorized(withRole WithRole, askedPerm Permission) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	role := withRole.Role()

	if role == nil {
		return false
	}

	permissions, ok := s.permsByRole[*role]
	if !ok {
		return false
	}

	for _, perm := range permissions {
		if perm == askedPerm {
			return true
		}
	}

	return false
}

func (s *service) createDefaultRoles(ctx context.Context) error {
	for role, permissions := range DefaultRoles {
		err := s.storage.Save(ctx, &role, permissions)
		if err != nil {
			return fmt.Errorf("failed to save the role %q: %w", role, err)
		}

		s.logger.Info(fmt.Sprintf("role %q created", role))
	}

	return nil
}
