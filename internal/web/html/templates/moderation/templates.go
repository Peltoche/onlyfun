package moderation

import (
	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/partials"
)

type OverviewPageTmpl struct {
	Header                 *partials.HeaderTmpl
	PostsWaitingModeration int
}

func (t *OverviewPageTmpl) Template() string { return "moderation/page_overview" }

type NextPostsPageTmpl struct {
	Header       *partials.HeaderTmpl
	Post         *posts.Post
	Media        *medias.FileMeta
	Author       *users.User
	AuthorAvatar *medias.FileMeta
	AuthorStats  map[posts.Status]int
}

func (t *NextPostsPageTmpl) Template() string { return "moderation/page_next_post" }
