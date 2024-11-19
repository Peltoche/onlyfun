package partials

import "github.com/Peltoche/onlyfun/internal/services/users"

type HeaderTmpl struct {
	User        *users.User
	CanModerate bool
	PostButton  bool
}
