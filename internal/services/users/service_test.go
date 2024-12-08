package users

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Users_Service(t *testing.T) {
	ctx := context.Background()

	t.Run("Create success", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		mediasSvc := medias.NewMockService(t)
		services := newService(tools, storage, mediasSvc)

		// Data
		role, _ := perms.NewFakePermissions(t).Build()
		avatar := medias.NewFakeFileMeta(t).Build()
		user := NewFakeUser(t).Build()
		newUser := NewFakeUser(t).
			CreatedBy(user).
			WithAvatar(avatar).
			WithRole(role).
			Build()

		// Mocks
		storage.On("GetByUsername", ctx, newUser.username).Return(nil, errNotFound).Once()

		mediasSvc.On("Upload", mock.Anything, medias.Avatar, mock.Anything).Return(avatar, nil).Once()

		tools.UUIDMock.On("New").Return(newUser.id).Once()

		tools.ClockMock.On("Now").Return(newUser.createdAt).Once()
		tools.PasswordMock.On("Encrypt", ctx, secret.NewText("my-super-password")).
			Return(newUser.password, nil).Once()

		storage.On("Save", ctx, newUser).Return(nil)

		// Run
		res, err := services.Create(ctx, &CreateCmd{
			CreatedBy: user,
			Role:      role,
			Username:  newUser.username,
			Password:  secret.NewText("my-super-password"),
		})

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, newUser, res)
	})

	t.Run("Create with a taken username", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		role, _ := perms.NewFakePermissions(t).Build()
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetByUsername", ctx, "Donald-Duck").Return(&User{}, nil).Once()

		// Run
		res, err := services.Create(ctx, &CreateCmd{
			CreatedBy: user,
			Role:      role,
			Username:  "Donald-Duck",
			Password:  secret.NewText("some-password"),
		})

		// Asserts
		require.ErrorIs(t, err, ErrUsernameTaken)
		require.ErrorIs(t, err, errs.ErrBadRequest)
		assert.Nil(t, res)
	})

	t.Run("Create with a database error", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		role, _ := perms.NewFakePermissions(t).Build()
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetByUsername", ctx, "Donald-Duck").Return(nil, fmt.Errorf("some-error")).Once()

		res, err := services.Create(ctx, &CreateCmd{
			CreatedBy: user,
			Role:      role,
			Username:  "Donald-Duck",
			Password:  secret.NewText("some-secret"),
		})

		// Asserts
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		assert.Nil(t, res)
	})

	t.Run("Authenticate success", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetByUsername", ctx, "Donald-Duck").Return(user, nil).Once()
		tools.PasswordMock.On("Compare", ctx, user.password, secret.NewText("some-password")).Return(true, nil).Once()

		// Run
		res, err := services.Authenticate(ctx, "Donald-Duck", secret.NewText("some-password"))

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, user, res)
	})

	t.Run("Authenticate with an invalid username", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data

		// Mocks
		storage.On("GetByUsername", ctx, "Donald-Duck").Return(nil, errNotFound).Once()

		// Run
		res, err := services.Authenticate(ctx, "Donald-Duck", secret.NewText("some-secret"))

		// Asserts
		require.ErrorIs(t, err, errs.ErrBadRequest)
		require.ErrorIs(t, err, ErrInvalidUsername)
		assert.Nil(t, res)
	})

	t.Run("Authenticate with an invalid password", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetByUsername", ctx, "Donald-Duck").Return(user, nil).Once()
		tools.PasswordMock.On("Compare", ctx, user.password, secret.NewText("some-invalid-password")).Return(false, nil).Once()

		// Invalid password here
		res, err := services.Authenticate(ctx, "Donald-Duck", secret.NewText("some-invalid-password"))
		require.ErrorIs(t, err, ErrInvalidPassword)
		assert.Nil(t, res)
	})

	t.Run("Authenticate an unhandled password error", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetByUsername", ctx, "Donald-Duck").Return(user, nil).Once()
		tools.PasswordMock.On("Compare", ctx, user.password, secret.NewText("some-password")).Return(false, fmt.Errorf("some-error")).Once()

		// Run
		res, err := services.Authenticate(ctx, "Donald-Duck", secret.NewText("some-password"))

		// Asserts
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
		assert.Nil(t, res)
	})

	t.Run("GetByID success", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetByID", ctx, user.ID()).Return(user, nil).Once()

		// Run
		res, err := services.GetByID(ctx, user.ID())

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, user, res)
	})

	t.Run("GetAll success", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetAll", ctx, &sqlstorage.PaginateCmd{Limit: 10}).Return([]User{*user}, nil).Once()

		// Run
		res, err := services.GetAll(ctx, &sqlstorage.PaginateCmd{Limit: 10})

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, []User{*user}, res)
	})

	t.Run("GetAllWithStatus success", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetAll", ctx, &sqlstorage.PaginateCmd{Limit: 10}).Return([]User{*user}, nil).Once()

		// Run
		res, err := services.GetAllWithStatus(ctx, Active, &sqlstorage.PaginateCmd{Limit: 10})

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, []User{*user}, res)
	})

	// t.Run("AddToDeletion an admin user success", func(t *testing.T) {
	// 	t.Parallel()
	// 	tools := tools.NewMock(t)
	// 	store := newMockStorage(t)
	// 	services := newService(tools, store)
	//
	// 	// Data
	// 	user := NewFakeUser(t).WithAdminRole().Build()
	// 	anAnotherAdmin := NewFakeUser(t).WithAdminRole().Build()
	//
	// 	// Mocks
	// 	store.On("GetByID", ctx, user.ID()).Return(user, nil).Once()
	// 	store.On("GetAll", ctx, (*sqlstorage.PaginateCmd)(nil)).
	// 		Return([]User{*user, *anAnotherAdmin}, nil).Once() // We check that the deleted user is not the last admin.
	// 	store.On("HardDelete", mock.Anything, user.ID()).Return(nil).Once()
	//
	// 	// Run
	// 	err := services.AddToDeletion(ctx, user.ID())
	//
	// 	// Asserts
	// 	require.NoError(t, err)
	// })

	t.Run("AddToDeletion with a user not found", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetByID", ctx, user.ID()).Return(nil, errNotFound).Once()

		// Run
		err := services.AddToDeletion(ctx, user.ID())

		// Asserts
		require.ErrorIs(t, err, errs.ErrNotFound)
		require.ErrorIs(t, err, errNotFound)
	})

	// t.Run("AddToDeletion the last admin failed", func(t *testing.T) {
	// 	t.Parallel()
	// 	tools := tools.NewMock(t)
	// 	store := newMockStorage(t)
	// 	services := newService(tools, store)
	//
	// 	// Data
	// 	user := NewFakeUser(t).WithAdminRole().Build()
	// 	anAnotherUser := NewFakeUser(t).Build()
	//
	// 	// Mocks
	// 	store.On("GetByID", ctx, user.ID()).Return(user, nil).Once()
	// 	store.On("GetAll", ctx, (*sqlstorage.PaginateCmd)(nil)).Return([]User{*user, *anAnotherUser}, nil).Once() // This is the last admin
	//
	// 	err := services.AddToDeletion(ctx, user.ID())
	// 	require.EqualError(t, err, "unauthorized: can't remove the last admin")
	// })
	//
	t.Run("HardDelete success", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		someSoftDeletedUser := NewFakeUser(t).WithStatus(Deleting).Build()

		// Mocks
		storage.On("GetByID", mock.Anything, someSoftDeletedUser.ID()).Return(someSoftDeletedUser, nil).Once()
		storage.On("HardDelete", mock.Anything, someSoftDeletedUser.ID()).Return(nil).Once()

		// Run
		err := services.HardDelete(ctx, someSoftDeletedUser.ID())

		// Asserts
		require.NoError(t, err)
	})

	t.Run("HardDelete an non existing user", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		someSoftDeletedUser := NewFakeUser(t).WithStatus(Deleting).Build()

		// Mocks
		storage.On("GetByID", mock.Anything, someSoftDeletedUser.ID()).Return(nil, errNotFound).Once()

		// Run
		err := services.HardDelete(ctx, someSoftDeletedUser.ID())

		// Asserts
		require.NoError(t, err)
	})

	t.Run("HardDelete an invalid status", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		someStillActifUser := NewFakeUser(t).WithStatus(Active).Build()

		// Mocks
		storage.On("GetByID", mock.Anything, someStillActifUser.ID()).Return(someStillActifUser, nil).Once()

		// Run
		err := services.HardDelete(ctx, someStillActifUser.ID())

		// Asserts
		require.ErrorIs(t, err, ErrInvalidStatus)
	})

	t.Run("UpdatePassword success", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()
		now := time.Now()

		// Mocks
		storage.On("GetByID", mock.Anything, user.ID()).Return(user, nil).Once()
		tools.PasswordMock.On("Encrypt", mock.Anything, secret.NewText("some-new-password")).
			Return(secret.NewText("some-encrypted-password"), nil).Once()
		tools.ClockMock.On("Now").Return(now).Once()
		storage.On("Patch", mock.Anything, user.ID(), map[string]any{
			"password":            secret.NewText("some-encrypted-password"),
			"password_changed_at": sqlstorage.SQLTime(now),
		}).Return(nil).Once()

		// Run
		err := services.UpdateUserPassword(ctx, &UpdatePasswordCmd{
			UserID:      user.ID(),
			NewPassword: secret.NewText("some-new-password"),
		})

		// Asserts
		require.NoError(t, err)
	})

	t.Run("UpdatePassword with a user not found", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()

		// Mocks
		storage.On("GetByID", mock.Anything, user.ID()).
			Return(nil, errs.ErrNotFound).Once()

		// Run
		err := services.UpdateUserPassword(ctx, &UpdatePasswordCmd{
			UserID:      user.ID(),
			NewPassword: secret.NewText("some-password"),
		})

		// Asserts
		require.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("UpdatePassword with a patch error", func(t *testing.T) {
		t.Parallel()
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		medias := medias.NewMockService(t)
		services := newService(tools, storage, medias)

		// Data
		user := NewFakeUser(t).Build()
		now := time.Now()

		// Mocks
		storage.On("GetByID", mock.Anything, user.ID()).Return(user, nil).Once()

		tools.PasswordMock.On("Encrypt", mock.Anything, secret.NewText("some-new-password")).
			Return(secret.NewText("some-encrypted-password"), nil).Once()

		tools.ClockMock.On("Now").Return(now).Once()

		storage.On("Patch", mock.Anything, user.ID(), map[string]any{
			"password":            secret.NewText("some-encrypted-password"),
			"password_changed_at": sqlstorage.SQLTime(now),
		}).Return(fmt.Errorf("some-error")).Once()

		// Run
		err := services.UpdateUserPassword(ctx, &UpdatePasswordCmd{
			UserID:      user.ID(),
			NewPassword: secret.NewText("some-new-password"),
		})

		// Asserts
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
	})
}
