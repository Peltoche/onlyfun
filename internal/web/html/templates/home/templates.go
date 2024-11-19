package home

import (
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/partials"
)

type ListingPageTmpl struct {
	Header *partials.HeaderTmpl
	Posts  []posts.Post
}

func (t *ListingPageTmpl) Template() string { return "home/page_listing" }

type SubmitPageTmpl struct {
	Header *partials.HeaderTmpl
}

func (t *SubmitPageTmpl) Template() string { return "home/page_submit" }
