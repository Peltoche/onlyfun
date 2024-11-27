package roles

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
	Save(ctx context.Context, r *Role) error
	GetAll(ctx context.Context) ([]Role, error)
	GetByName(ctx context.Context, roleName string) (*Role, error)
}

type service struct {
	storage storage
	clock   clock.Clock
	uuid    uuid.Service
	logger  *slog.Logger

	roles map[string]Role
	lock  *sync.RWMutex
}

func newService(tools tools.Tools, storage storage) *service {
	return &service{
		storage: storage,
		clock:   tools.Clock(),
		uuid:    tools.UUID(),
		logger:  tools.Logger(),

		roles: map[string]Role{},
		lock:  new(sync.RWMutex),
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

	for _, r := range roles {
		s.roles[r.name] = r
	}

	return nil
}

func (s *service) IsRoleAuthorized(roleName string, askedPerm Permission) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	role, ok := s.roles[roleName]
	if !ok {
		return false
	}

	for _, perm := range role.permissions {
		if perm == askedPerm {
			return true
		}
	}

	return false
}

func (s *service) createDefaultRoles(ctx context.Context) error {
	err := s.storage.Save(ctx, &DefaultAdminRole)
	if err != nil {
		return fmt.Errorf("failed to save the DefaultAdminRole: %w", err)
	}
	s.logger.Info(fmt.Sprintf("Default role %q created", DefaultAdminRole.name))

	err = s.storage.Save(ctx, &DefaultModeratorRole)
	if err != nil {
		return fmt.Errorf("failed to save the DefaultModeratorRole: %w", err)
	}
	s.logger.Info(fmt.Sprintf("Default role %q created", DefaultModeratorRole.name))

	err = s.storage.Save(ctx, &DefaultUserRole)
	if err != nil {
		return fmt.Errorf("failed to save the DefaultUserRole: %w", err)
	}
	s.logger.Info(fmt.Sprintf("Default role %q created", DefaultUserRole.name))

	return nil
}
