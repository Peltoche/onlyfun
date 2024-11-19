package moderation

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/roles"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/services/websessions"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/router"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/Peltoche/onlyfun/internal/web/handlers/auth"
	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/moderation"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/partials"
	"github.com/go-chi/chi/v5"
)

type ModerationHandler struct {
	webSessions websessions.Service
	auth        *auth.Authenticator
	posts       posts.Service
	medias      medias.Service
	users       users.Service
	roles       roles.Service
	html        html.Writer
	uuid        uuid.Service
}

func NewModerationHandler(
	ctx context.Context,
	html html.Writer,
	auth *auth.Authenticator,
	posts posts.Service,
	users users.Service,
	roles roles.Service,
	medias medias.Service,
	tools tools.Tools,
) *ModerationHandler {
	return &ModerationHandler{
		html:   html,
		posts:  posts,
		users:  users,
		medias: medias,
		roles:  roles,
		auth:   auth,
	}
}

func (h *ModerationHandler) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/moderation", h.printOverviewPage)
	r.Get("/moderation/posts", h.printPostsPage)
}

func (h *ModerationHandler) printOverviewPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, _, err := h.auth.GetUserAndSession(w, r)
	if err != nil && !errors.Is(err, auth.ErrNotAuthenticated) {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	if user == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if !h.roles.IsRoleAuthorized(user.Role(), roles.Moderation) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	waitingModeration, err := h.posts.CountPostsWaitingModeration(ctx)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to CountPostsWaitingModeration: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &moderation.OverviewPageTmpl{
		Header: &partials.HeaderTmpl{
			User:        user,
			CanModerate: h.roles.IsRoleAuthorized(user.Role(), roles.Moderation),
			PostButton:  false,
		},

		PostsWaitingModeration: waitingModeration,
	})
}

func (h *ModerationHandler) printPostsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, _, err := h.auth.GetUserAndSession(w, r)
	if err != nil && !errors.Is(err, auth.ErrNotAuthenticated) {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	if user == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if !h.roles.IsRoleAuthorized(user.Role(), roles.Moderation) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	post, err := h.posts.GetNextPostToModerate(ctx)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to GetNextPostToMederate"))
		return
	}

	fileMeta, err := h.medias.GetMetadata(ctx, post.FileID())
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	author, err := h.users.GetByID(ctx, post.CreatedBy())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the author: %w", err))
		return
	}

	avatarMeta, err := h.medias.GetMetadata(ctx, author.Avatar())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the author avatar: %w", err))
		return
	}

	stats, err := h.posts.GetUserStats(ctx, user)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the author posts stats: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &moderation.NextPostsPageTmpl{
		Header: &partials.HeaderTmpl{
			User:        user,
			CanModerate: h.roles.IsRoleAuthorized(user.Role(), roles.Moderation),
			PostButton:  false,
		},

		Post:         post,
		Media:        fileMeta,
		Author:       author,
		AuthorAvatar: avatarMeta,
		AuthorStats:  stats,
	})
}
