package home

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/router"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/Peltoche/onlyfun/internal/web/handlers/auth"
	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/home"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/partials"
	"github.com/go-chi/chi/v5"
)

const postPagination = 50

type ListingPage struct {
	roles      perms.Service
	auth       *auth.Authenticator
	posts      posts.Service
	medias     medias.Service
	html       html.Writer
	l          *sync.Mutex
	latestPost *posts.Post
	uuid       uuid.Service
}

func NewListingPage(
	ctx context.Context,
	html html.Writer,
	posts posts.Service,
	roles perms.Service,
	auth *auth.Authenticator,
	medias medias.Service,
	tools tools.Tools,
) (*ListingPage, error) {
	latest, err := posts.GetLatestPost(ctx)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return nil, fmt.Errorf("failed to GetLatestPost: %w", err)
	}

	handler := &ListingPage{
		html:       html,
		uuid:       uuid.NewProvider(),
		posts:      posts,
		roles:      roles,
		medias:     medias,
		auth:       auth,
		l:          new(sync.Mutex),
		latestPost: latest,
	}

	postChan := posts.SuscribeToNewPost()

	go func() {
		for post := range postChan {
			p := post
			handler.l.Lock()
			handler.latestPost = &p
			handler.l.Unlock()
		}
	}()

	return handler, nil
}

func (h *ListingPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/", h.printPage)
	r.Get("/medias/{fileID}", h.serveMedia)
}

func (h *ListingPage) printPage(w http.ResponseWriter, r *http.Request) {
	var err error
	var posts []posts.Post

	user, _, err := h.auth.GetUserAndSession(w, r)
	if err != nil && !errors.Is(err, auth.ErrNotAuthenticated) {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	if h.latestPost != nil {
		posts, err = h.posts.GetPosts(r.Context(), h.latestPost.ID(), postPagination)
		if err != nil {
			h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to GetPosts: %w", err))
			return
		}
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &home.ListingPageTmpl{
		Header: &partials.HeaderTmpl{
			User:        user,
			CanModerate: user != nil && h.roles.IsAuthorized(user, perms.Moderation),
			PostButton:  true,
		},
		Posts: posts,
	})
}

func (h *ListingPage) serveMedia(w http.ResponseWriter, r *http.Request) {
	fileID, err := h.uuid.Parse(chi.URLParam(r, "fileID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.medias.GetMetadata(r.Context(), fileID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := h.medias.Download(r.Context(), fileID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res != nil {
		w.Header().Set("ETag", fmt.Sprintf("W/%q", res.Checksum()))
		w.Header().Set("Content-Type", res.Mimetype())
	}

	w.Header().Set("Expires", time.Now().Add(365*24*time.Hour).UTC().Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "max-age=31536000")

	http.ServeContent(w, r, string(res.ID()), res.UploadedAt(), content)
}
