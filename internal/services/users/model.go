package users

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/roles"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	SecretMinLength = 8
	SecretMaxLength = 200
)

var UsernameRegexp = regexp.MustCompile("^[0-9a-zA-Z-]+$")

type Status string

const (
	Active   Status = "active"
	Deleting Status = "deleting"
)

// User representation
type User struct {
	createdAt         time.Time
	passwordChangedAt time.Time
	id                uuid.UUID
	username          string
	password          secret.Text
	role              string
	status            Status
	avatar            uuid.UUID // Media's id
	createdBy         uuid.UUID
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":        u.id,
		"username":  u.username,
		"role":      u.role,
		"createdAt": u.createdAt,
		"avatar":    u.avatar,
		"status":    u.status,
	})
}

func (u User) ID() uuid.UUID                { return u.id }
func (u User) Username() string             { return u.username }
func (u User) Role() string                 { return u.role }
func (u User) Status() Status               { return u.status }
func (u User) PasswordChangedAt() time.Time { return u.passwordChangedAt }
func (u User) Avatar() uuid.UUID            { return u.avatar }
func (u User) CreatedAt() time.Time         { return u.createdAt }
func (u User) CreatedBy() uuid.UUID         { return u.createdBy }

// CreateCmd represents an user creation request.
type CreateCmd struct {
	CreatedBy *User
	Role      *roles.Role
	Username  string
	Password  secret.Text
}

// Validate the CreateUserRequest fields.
func (t CreateCmd) Validate() error {
	return v.ValidateStruct(&t,
		v.Field(&t.CreatedBy, v.Required),
		v.Field(&t.Role, v.Required),
		v.Field(&t.Username, v.Required, v.Length(1, 20), v.Match(UsernameRegexp)),
		v.Field(&t.Password, v.Required, v.Length(SecretMinLength, SecretMaxLength)),
	)
}

type UpdatePasswordCmd struct {
	UserID      uuid.UUID
	NewPassword secret.Text
}

func (t UpdatePasswordCmd) Validate() error {
	return v.ValidateStruct(&t,
		v.Field(&t.UserID, v.Required, is.UUIDv4),
		v.Field(&t.NewPassword, v.Required, v.Length(SecretMinLength, SecretMaxLength)),
	)
}

type BootstrapCmd struct {
	Username string
	Password secret.Text
}

func (t BootstrapCmd) Validate() error {
	return v.ValidateStruct(&t,
		v.Field(&t.Username, v.Required, v.Length(1, 20), v.Match(UsernameRegexp)),
		v.Field(&t.Password, v.Required, v.Length(SecretMinLength, SecretMaxLength)),
	)
}
