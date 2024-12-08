package users

import (
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_User_Getters(t *testing.T) {
	u := NewFakeUser(t).Build()

	assert.Equal(t, u.id, u.ID())
	assert.Equal(t, u.role, u.Role())
	assert.Equal(t, u.username, u.Username())
	assert.Equal(t, u.createdAt, u.CreatedAt())
	assert.Equal(t, u.avatar, u.Avatar())
	assert.Equal(t, u.passwordChangedAt, u.PasswordChangedAt())
	assert.Equal(t, u.createdBy, u.CreatedBy())
	assert.Equal(t, u.status, u.Status())
}

func Test_CreateCmd_is_validatable(t *testing.T) {
	assert.Implements(t, (*v.Validatable)(nil), new(CreateCmd))
}

func Test_CreateCmd_Validate_success(t *testing.T) {
	role, _ := perms.NewFakePermissions(t).Build()

	err := CreateCmd{
		CreatedBy: NewFakeUser(t).Build(),
		Role:      role,
		Username:  "some-username",
		Password:  secret.NewText("myLittleSecret"),
	}.Validate()

	require.NoError(t, err)
}

func Test_UpdatePasswordCmd(t *testing.T) {
	err := UpdatePasswordCmd{
		UserID:      uuid.UUID("some-invalid-id"),
		NewPassword: secret.NewText("foobar1234"),
	}.Validate()

	require.EqualError(t, err, "UserID: must be a valid UUID v4.")
}

func Test_BootstrapCmd_is_validatable(t *testing.T) {
	assert.Implements(t, (*v.Validatable)(nil), new(BootstrapCmd))
}

func Test_BootstrapCmd_Validate_success(t *testing.T) {
	err := BootstrapCmd{
		Username: "some-username",
		Password: secret.NewText("myLittleSecret"),
	}.Validate()

	require.NoError(t, err)
}
