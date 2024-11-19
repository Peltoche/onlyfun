package users

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/png"

	"github.com/o1egl/govatar"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/roles"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/password"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

var (
	ErrAlreadyExists     = fmt.Errorf("user already exists")
	ErrUsernameTaken     = fmt.Errorf("username taken")
	ErrInvalidUsername   = fmt.Errorf("invalid username")
	ErrInvalidPassword   = fmt.Errorf("invalid password")
	ErrLastAdmin         = fmt.Errorf("can't remove the last admin")
	ErrInvalidStatus     = fmt.Errorf("invalid status")
	ErrUnauthorizedSpace = fmt.Errorf("unauthorized space")
)

// storage encapsulates the logic to access user from the data source.
type storage interface {
	Save(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetAll(ctx context.Context, cmd *sqlstorage.PaginateCmd) ([]User, error)
	HardDelete(ctx context.Context, userID uuid.UUID) error
	Patch(ctx context.Context, userID uuid.UUID, fields map[string]any) error
}

// services handling all the logic.
type services struct {
	medias   medias.Service
	storage  storage
	clock    clock.Clock
	uuid     uuid.Service
	password password.Password
}

// newService create a new user services.
func newService(tools tools.Tools, storage storage, medias medias.Service) *services {
	return &services{
		medias:   medias,
		storage:  storage,
		clock:    tools.Clock(),
		uuid:     tools.UUID(),
		password: tools.Password(),
	}
}

func (s *services) Bootstrap(ctx context.Context, cmd *BootstrapCmd) (*User, error) {
	err := cmd.Validate()
	if err != nil {
		return nil, errs.Validation(err)
	}

	newUserID := s.uuid.New()
	return s.createUser(ctx, newUserID, &roles.DefaultAdminRole, cmd.Username, cmd.Password, newUserID)
}

// Create will create and register a new user.
func (s *services) Create(ctx context.Context, cmd *CreateCmd) (*User, error) {
	err := cmd.Validate()
	if err != nil {
		return nil, errs.Validation(err)
	}

	userWithSameUsername, err := s.storage.GetByUsername(ctx, cmd.Username)
	if err != nil && !errors.Is(err, errNotFound) {
		return nil, errs.Internal(fmt.Errorf("failed to GetByUsername: %w", err))
	}

	if userWithSameUsername != nil {
		return nil, errs.BadRequest(ErrUsernameTaken, "username already taken")
	}

	newUserID := s.uuid.New()
	return s.createUser(ctx, newUserID, cmd.Role, cmd.Username, cmd.Password, cmd.CreatedBy.id)
}

func (s *services) createUser(
	ctx context.Context,
	newUserID uuid.UUID,
	role *roles.Role,
	username string,
	password secret.Text,
	createdBy uuid.UUID) (*User, error) {
	hashedPassword, err := s.password.Encrypt(ctx, password)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to hash the password: %w", err))
	}

	now := s.clock.Now()

	avatarImg, err := govatar.GenerateForUsername(govatar.MALE, username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate the avatar: %w", err)
	}

	var avatarBuf bytes.Buffer
	err = png.Encode(&avatarBuf, avatarImg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode the image into png: %w", err)
	}

	avatar, err := s.medias.Upload(ctx, medias.Avatar, &avatarBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to save the avatar: %w", err)
	}

	user := User{
		id:                newUserID,
		role:              role.Name(),
		username:          username,
		password:          hashedPassword,
		status:            Active,
		passwordChangedAt: now,
		avatar:            avatar.ID(),
		createdAt:         now,
		createdBy:         createdBy,
	}

	err = s.storage.Save(ctx, &user)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to save the user: %w", err))
	}

	return &user, nil
}

func (s *services) UpdateUserPassword(ctx context.Context, cmd *UpdatePasswordCmd) error {
	user, err := s.GetByID(ctx, cmd.UserID)
	if err != nil {
		return fmt.Errorf("failed to GetByID: %w", err)
	}

	hashedPassword, err := s.password.Encrypt(ctx, cmd.NewPassword)
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to hash the password: %w", err))
	}

	err = s.storage.Patch(ctx, user.ID(), map[string]any{
		"password":            hashedPassword,
		"password_changed_at": sqlstorage.SQLTime(s.clock.Now()),
	})
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to patch the user: %w", err))
	}

	return nil
}

func (s *services) GetAllWithStatus(ctx context.Context, status Status, cmd *sqlstorage.PaginateCmd) ([]User, error) {
	allUsers, err := s.GetAll(ctx, cmd)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to GetAll users: %w", err))
	}

	res := []User{}
	for _, user := range allUsers {
		if user.status == status {
			res = append(res, user)
		}
	}

	return res, nil
}

// Authenticate return the user corresponding to the given username only if the password is correct.
func (s *services) Authenticate(ctx context.Context, username string, userPassword secret.Text) (*User, error) {
	user, err := s.storage.GetByUsername(ctx, username)
	if errors.Is(err, errNotFound) {
		return nil, errs.BadRequest(ErrInvalidUsername)
	}
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to GetbyUsername: %w", err))
	}

	ok, err := s.password.Compare(ctx, user.password, userPassword)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed password compare: %w", err))
	}

	if !ok {
		return nil, errs.BadRequest(ErrInvalidPassword)
	}

	return user, nil
}

func (s *services) GetByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	res, err := s.storage.GetByID(ctx, userID)
	if errors.Is(err, errNotFound) {
		return nil, errs.NotFound(err)
	}

	if err != nil {
		return nil, errs.Internal(err)
	}

	return res, nil
}

func (s *services) GetAll(ctx context.Context, paginateCmd *sqlstorage.PaginateCmd) ([]User, error) {
	res, err := s.storage.GetAll(ctx, paginateCmd)
	if err != nil {
		return nil, errs.Internal(err)
	}

	return res, nil
}

func (s *services) AddToDeletion(ctx context.Context, userID uuid.UUID) error {
	_, err := s.GetByID(ctx, userID)
	if errors.Is(err, errNotFound) {
		return errs.NotFound(err)
	}

	if err != nil {
		return errs.Internal(fmt.Errorf("failed to GetByID: %w", err))
	}

	// TODO: Should check if the user is not the last one with the user creation
	// if user.IsAdmin() {
	// 	users, err := s.GetAll(ctx, nil)
	// 	if err != nil {
	// 		return errs.Internal(fmt.Errorf("failed to GetAll: %w", err))
	// 	}
	//
	// 	if isTheLastAdmin(users) {
	// 		return errs.Unauthorized(ErrLastAdmin, "you are the last admin, you account can't be removed")
	// 	}
	// }

	err = s.storage.HardDelete(ctx, userID)
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to Patch the user: %w", err))
	}

	return nil
}

func (s *services) HardDelete(ctx context.Context, userID uuid.UUID) error {
	res, err := s.storage.GetByID(ctx, userID)
	if errors.Is(err, errNotFound) {
		return nil
	}
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to GetDeleted: %w", err))
	}

	if res.status != Deleting {
		return errs.Internal(ErrInvalidStatus)
	}

	err = s.storage.HardDelete(ctx, userID)
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to HardDelete: %w", err))
	}

	return nil
}
