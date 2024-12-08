package perms

type Permission string

const (
	UploadPost Permission = "posts.upload"
	Moderation Permission = "moderation"
)

type Role string

const (
	DefaultAdminRole     Role = "admin"
	DefaultModeratorRole Role = "moderator"
	DefaultUserRole      Role = "user"
)

var DefaultRoles = map[Role][]Permission{
	DefaultAdminRole:     {UploadPost, Moderation},
	DefaultModeratorRole: {UploadPost, Moderation},
	DefaultUserRole:      {UploadPost},
}
