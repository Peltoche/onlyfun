package home

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/router"
	"github.com/Peltoche/onlyfun/internal/web/handlers/auth"
	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/home"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/partials"
	"github.com/go-chi/chi/v5"
)

type SubmitPage struct {
	posts posts.Service
	roles perms.Service
	auth  *auth.Authenticator
	html  html.Writer
}

func NewSubmitPage(
	html html.Writer,
	auth *auth.Authenticator,
	posts posts.Service,
	roles perms.Service,
	tools tools.Tools,
) *SubmitPage {
	return &SubmitPage{
		html:  html,
		posts: posts,
		roles: roles,
		auth:  auth,
	}
}

func (h *SubmitPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/submit", h.printPage)
	r.Post("/submit", h.submitPost)
}

func (h *SubmitPage) printPage(w http.ResponseWriter, r *http.Request) {
	user, _, err := h.auth.GetUserAndSession(w, r)
	if err != nil && !errors.Is(err, auth.ErrNotAuthenticated) {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	if errors.Is(err, auth.ErrNotAuthenticated) {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &home.SubmitPageTmpl{
		Header: &partials.HeaderTmpl{
			User:        user,
			CanModerate: h.roles.IsAuthorized(user, perms.Moderation),
			PostButton:  true,
		},
	})
}

func (h *SubmitPage) submitPost(w http.ResponseWriter, r *http.Request) {
	user, _, err := h.auth.GetUserAndSession(w, r)
	if err != nil && !errors.Is(err, auth.ErrNotAuthenticated) {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to ParseForm: %w", err))
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to retrieve the FormFile: %w", err))
		return
	}
	defer file.Close()

	_, err = h.posts.Create(r.Context(), &posts.CreateCmd{
		Title:     r.FormValue("title"),
		Media:     file,
		CreatedBy: user,
	})
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to create the post: %w", err))
		return
	}

	http.Redirect(w, r, "", http.StatusFound)
}
