package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/services/websessions"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/web/html"
)

var ErrNotAuthenticated = errors.New("user not authenticated")

type AccessType int

const (
	AdminOnly AccessType = iota
	AnyUser
)

type Authenticator struct {
	webSessions websessions.Service
	users       users.Service
	html        html.Writer
}

func NewAuthenticator(webSessions websessions.Service, users users.Service, html html.Writer) *Authenticator {
	return &Authenticator{webSessions, users, html}
}

func (a *Authenticator) GetUserAndSession(w http.ResponseWriter, r *http.Request) (*users.User, *websessions.Session, error) {
	currentSession, err := a.webSessions.GetFromReq(r)
	switch {
	case err == nil:
		break
	case errors.Is(err, websessions.ErrSessionNotFound):
		a.webSessions.Logout(r, w)
		return nil, nil, ErrNotAuthenticated
	case errors.Is(err, websessions.ErrMissingSessionToken):
		return nil, nil, ErrNotAuthenticated
	default:
		return nil, nil, errs.Internal(fmt.Errorf("failed to websessions.GetFromReq: %w", err))
	}

	user, err := a.users.GetByID(r.Context(), currentSession.UserID())
	if err != nil {
		return nil, nil, errs.Internal(fmt.Errorf("failed to users.GetByID: %w", err))
	}

	if user == nil {
		_ = a.webSessions.Logout(r, w)
		return nil, nil, ErrNotAuthenticated
	}

	return user, currentSession, nil
}

func (a *Authenticator) Logout(w http.ResponseWriter, r *http.Request) {
	a.webSessions.Logout(r, w)
}
