package roles

import (
	"context"
	"fmt"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/stretchr/testify/require"
)

func Test_Roles_Service(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("bootstrap and IsRoleAuthorized success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage)

		userRole := NewFakeRole(t).WithName("user").WithPermissions(UploadPost).Build()
		adminRole := NewFakeRole(t).WithName("admin").WithPermissions(UploadPost, Moderation).Build()

		storage.On("GetAll", ctx).Return([]Role{*userRole, *adminRole}, nil).Once()

		err := svc.bootstrap(ctx)
		require.NoError(t, err)

		res1 := svc.IsRoleAuthorized("user", UploadPost)
		require.True(t, res1) // The user can upload posts.

		res2 := svc.IsRoleAuthorized("user", Moderation)
		require.False(t, res2) // The user don't have moderation write
	})

	t.Run("IsRoleAuthorized with an invalid role", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage)

		res := svc.IsRoleAuthorized("unknown-role", UploadPost)
		require.False(t, res)
	})

	t.Run("boostrap and createDefaultRoles success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage)

		storage.On("GetAll", ctx).Return([]Role{}, nil).Once()
		storage.On("Save", ctx, &DefaultAdminRole).Return(nil).Once()
		storage.On("Save", ctx, &DefaultModeratorRole).Return(nil).Once()
		storage.On("Save", ctx, &DefaultUserRole).Return(nil).Once()
		storage.On("GetAll", ctx).Return([]Role{DefaultAdminRole, DefaultModeratorRole, DefaultUserRole}, nil).Once()

		err := svc.bootstrap(ctx)
		require.NoError(t, err)

		require.EqualValues(t, map[string]Role{
			DefaultAdminRole.name:     DefaultAdminRole,
			DefaultModeratorRole.name: DefaultModeratorRole,
			DefaultUserRole.name:      DefaultUserRole,
		}, svc.roles)
	})

	t.Run("createDefaultRoles with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage)

		storage.On("Save", ctx, &DefaultAdminRole).Return(nil).Once()
		storage.On("Save", ctx, &DefaultModeratorRole).Return(fmt.Errorf("some-error")).Once()

		err := svc.createDefaultRoles(ctx)
		require.ErrorContains(t, err, "some-error")
	})
}
