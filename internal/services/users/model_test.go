package users

import (
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/roles"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_User_Getters(t *testing.T) {
	u := NewFakeUser(t).Build()

	assert.Equal(t, u.id, u.ID())
	assert.Equal(t, u.role, u.Role())
	assert.Equal(t, u.username, u.Username())
	assert.Equal(t, u.createdAt, u.CreatedAt())
	assert.Equal(t, u.status, u.Status())
}

func Test_CreateUserRequest_is_validatable(t *testing.T) {
	assert.Implements(t, (*validation.Validatable)(nil), new(CreateCmd))
}

func Test_CreateUserRequest_Validate_success(t *testing.T) {
	err := CreateCmd{
		CreatedBy: NewFakeUser(t).Build(),
		Role:      roles.NewFakeRole(t).Build(),
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
