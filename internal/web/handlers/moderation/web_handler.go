package moderation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/moderations"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/router"
	"github.com/Peltoche/onlyfun/internal/web/handlers/auth"
	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/moderation"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/partials"
	"github.com/go-chi/chi/v5"
)

type ModerationHandler struct {
	auth      *auth.Authenticator
	postsSvc  posts.Service
	modeSvc   moderations.Service
	mediasSvc medias.Service
	usersSvc  users.Service
	permsSvc  perms.Service
	html      html.Writer
}

func NewModerationHandler(
	ctx context.Context,
	html html.Writer,
	auth *auth.Authenticator,
	posts posts.Service,
	modeSvc moderations.Service,
	users users.Service,
	roles perms.Service,
	medias medias.Service,
	tools tools.Tools,
) *ModerationHandler {
	return &ModerationHandler{
		html:      html,
		postsSvc:  posts,
		modeSvc:   modeSvc,
		usersSvc:  users,
		mediasSvc: medias,
		permsSvc:  roles,
		auth:      auth,
	}
}

func (h *ModerationHandler) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/moderation", h.printOverviewPage)
	r.Get("/moderation/posts", h.printPostsPage)
	r.Post("/moderation/posts", h.printPostsPage)
	r.Post("/moderation/posts/{postID}", h.handleValidation)
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

	if !h.permsSvc.IsAuthorized(user, perms.Moderation) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	waitingModeration, err := h.postsSvc.CountPostsWaitingModeration(ctx)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to CountPostsWaitingModeration: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &moderation.OverviewPageTmpl{
		Header: &partials.HeaderTmpl{
			User:        user,
			CanModerate: h.permsSvc.IsAuthorized(user, perms.Moderation),
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

	if !h.permsSvc.IsAuthorized(user, perms.Moderation) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	post, err := h.postsSvc.GetNextPostToModerate(ctx)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to GetNextPostToMederate"))
		return
	}

	fileMeta, err := h.mediasSvc.GetMetadata(ctx, post.FileID())
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	author, err := h.usersSvc.GetByID(ctx, post.CreatedBy())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the author: %w", err))
		return
	}

	avatarMeta, err := h.mediasSvc.GetMetadata(ctx, author.Avatar())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the author avatar: %w", err))
		return
	}

	stats, err := h.postsSvc.GetUserStats(ctx, user)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the author posts stats: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &moderation.NextPostsPageTmpl{
		Header: &partials.HeaderTmpl{
			User:        user,
			CanModerate: h.permsSvc.IsAuthorized(user, perms.Moderation),
			PostButton:  false,
		},

		Post:         post,
		Media:        fileMeta,
		Author:       author,
		AuthorAvatar: avatarMeta,
		AuthorStats:  stats,
	})
}

func (h *ModerationHandler) handleValidation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, _, err := h.auth.GetUserAndSession(w, r)
	if err != nil && !errors.Is(err, auth.ErrNotAuthenticated) {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	postID, err := strconv.ParseUint(chi.URLParam(r, "postID"), 10, 0)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to parse postID %q: %w", postID, err))
		return
	}

	post, err := h.postsSvc.GetByID(ctx, uint(postID))
	if errors.Is(err, errs.ErrNotFound) {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("post %q not found: %w", postID, err))
		return
	}

	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to GetByID: %w", err))
		return
	}

	isAccepted, err := strconv.ParseBool(r.FormValue("accepted"))
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("invalid value for accepted field %q: %w", r.FormValue("accepted"), err))
		return
	}

	switch isAccepted {
	case true:
		err = h.postsSvc.ValidatePost(ctx, &posts.ValidatePostcmd{
			User: user,
			Post: post,
		})
		if err != nil {
			h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to validate post %q: %w", postID, err))
			return
		}

	case false:
		_, err := h.modeSvc.ModeratePost(ctx, &moderations.PostModerationCmd{
			User:   user,
			Post:   post,
			Reason: r.FormValue("reason"),
		})
		if err != nil {
			h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to moderate post %q: %w", postID, err))
			return
		}
	}

	http.Redirect(w, r, "/moderation/posts", http.StatusTemporaryRedirect)
}
