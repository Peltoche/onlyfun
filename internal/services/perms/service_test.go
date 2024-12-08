package perms

import (
	"context"
	"fmt"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/ptr"
	"github.com/stretchr/testify/require"
)

type resourceWithRole struct {
	role *Role
}

func (r *resourceWithRole) Role() *Role {
	return r.role
}

func Test_Roles_Service(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("bootstrap and IsAuthorized success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage)

		userRole, userPerms := NewFakePermissions(t).WithName("user").WithPermissions(UploadPost).Build()
		adminRole, adminPerms := NewFakePermissions(t).WithName("admin").WithPermissions(UploadPost, Moderation).Build()

		storage.On("GetAll", ctx).Return(map[Role][]Permission{
			*userRole:  userPerms,
			*adminRole: adminPerms,
		}, nil).Once()

		resourceWithUserRole := resourceWithRole{userRole}

		err := svc.bootstrap(ctx)
		require.NoError(t, err)

		res1 := svc.IsAuthorized(&resourceWithUserRole, UploadPost)
		require.True(t, res1, "The user can upload posts.")

		res2 := svc.IsAuthorized(&resourceWithUserRole, Moderation)
		require.False(t, res2, "The user don't have moderation write")
	})

	t.Run("IsAuthorized with an invalid role", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage)

		resourceUnknownRole := resourceWithRole{ptr.To(Role("unknown-role"))}

		res := svc.IsAuthorized(&resourceUnknownRole, UploadPost)
		require.False(t, res)
	})

	t.Run("boostrap and createDefaultRoles success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage)

		storage.On("GetAll", ctx).Return(map[Role][]Permission{}, nil).Once()

		storage.On("Save", ctx, ptr.To(DefaultAdminRole), DefaultRoles[DefaultAdminRole]).Return(nil).Once()
		storage.On("Save", ctx, ptr.To(DefaultModeratorRole), DefaultRoles[DefaultModeratorRole]).Return(nil).Once()
		storage.On("Save", ctx, ptr.To(DefaultUserRole), DefaultRoles[DefaultUserRole]).Return(nil).Once()
		storage.On("GetAll", ctx).Return(DefaultRoles, nil).Once()

		err := svc.bootstrap(ctx)
		require.NoError(t, err)

		require.EqualValues(t, DefaultRoles, svc.permsByRole)
	})

	t.Run("createDefaultRoles with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage)

		storage.On("Save", ctx, ptr.To(DefaultAdminRole), DefaultRoles[DefaultAdminRole]).Return(nil).Once()
		storage.On("Save", ctx, ptr.To(DefaultModeratorRole), DefaultRoles[DefaultModeratorRole]).Return(fmt.Errorf("some-error")).Once()

		err := svc.createDefaultRoles(ctx)
		require.ErrorContains(t, err, "some-error")
	})
}
