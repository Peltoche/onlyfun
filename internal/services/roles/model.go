package roles

type Permission string

const (
	UploadPost Permission = "posts.upload"
	Moderation Permission = "moderation"
)

var (
	DefaultAdminRole = Role{
		name: "admin",
		permissions: []Permission{
			UploadPost,
			Moderation,
		},
	}
	DefaultModeratorRole = Role{
		name: "moderator",
		permissions: []Permission{
			UploadPost,
			Moderation,
		},
	}
	DefaultUserRole = Role{
		name: "user",
		permissions: []Permission{
			UploadPost,
		},
	}
)

type Role struct {
	name        string
	permissions []Permission
}

func (r Role) Name() string              { return r.name }
func (r Role) Permissions() []Permission { return r.permissions }
