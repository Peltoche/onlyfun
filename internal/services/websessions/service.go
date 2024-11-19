package websessions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	ua "github.com/mileusna/useragent"
)

var ErrUserIDNotMatching = errors.New("user ids are not matching")

type storage interface {
	Save(ctx context.Context, session *Session) error
	GetByToken(ctx context.Context, token secret.Text) (*Session, error)
	RemoveByToken(ctx context.Context, token secret.Text) error
	GetAllForUser(ctx context.Context, userID uuid.UUID, cmd *sqlstorage.PaginateCmd) ([]Session, error)
}

type services struct {
	clock   clock.Clock
	storage storage
	uuid    uuid.Service
}

func newService(storage storage, tools tools.Tools) *services {
	return &services{
		clock:   tools.Clock(),
		uuid:    tools.UUID(),
		storage: storage,
	}
}

func (s *services) Create(ctx context.Context, cmd *CreateCmd) (*Session, error) {
	err := cmd.Validate()
	if err != nil {
		return nil, errs.Validation(err)
	}

	uaRes := ua.Parse(cmd.UserAgent)

	session := &Session{
		token:     secret.NewText(string(s.uuid.New())),
		userID:    cmd.UserID,
		ip:        cmd.RemoteAddr,
		device:    fmt.Sprintf("%s - %s", uaRes.OS, uaRes.Name),
		createdAt: s.clock.Now(),
	}

	err = s.storage.Save(ctx, session)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to save the session: %w", err))
	}

	return session, nil
}

func (s *services) Delete(ctx context.Context, cmd *DeleteCmd) error {
	err := cmd.Validate()
	if err != nil {
		return errs.Validation(err)
	}

	session, err := s.storage.GetByToken(ctx, cmd.Token)
	if errors.Is(err, errNotFound) {
		return nil
	}

	if err != nil {
		return errs.Internal(fmt.Errorf("failed to GetByToken: %w", err))
	}

	if session.UserID() != cmd.UserID {
		return errs.NotFound(ErrUserIDNotMatching, "not found")
	}

	err = s.storage.RemoveByToken(ctx, session.Token())
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to RemoveByToken: %w", err))
	}

	return nil
}

func (s *services) GetByToken(ctx context.Context, token secret.Text) (*Session, error) {
	session, err := s.storage.GetByToken(ctx, token)
	if errors.Is(err, errNotFound) {
		return nil, errs.NotFound(err)
	}

	if err != nil {
		return nil, errs.Internal(err)
	}

	// TODO: Handle session expiration

	return session, nil
}

func (s *services) GetFromReq(r *http.Request) (*Session, error) {
	c, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		return nil, errs.BadRequest(ErrMissingSessionToken, "invalid_request")
	}

	session, err := s.GetByToken(r.Context(), secret.NewText(c.Value))
	if errors.Is(err, errNotFound) {
		return nil, errs.BadRequest(ErrSessionNotFound, "session not found")
	}

	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to GetByToken: %w", err))
	}

	return session, nil
}

func (s *services) Logout(r *http.Request, w http.ResponseWriter) error {
	c, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		// There is not session and so nothing to do.
		return nil
	}

	err = s.storage.RemoveByToken(r.Context(), secret.NewText(c.Value))
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to remove the token: %w", err))
	}

	// Remove to cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	return nil
}

func (s *services) GetAllForUser(ctx context.Context, userID uuid.UUID, cmd *sqlstorage.PaginateCmd) ([]Session, error) {
	res, err := s.storage.GetAllForUser(ctx, userID, cmd)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to GetAllForUser: %w", err))
	}

	return res, nil
}

func (s *services) DeleteAll(ctx context.Context, userID uuid.UUID) error {
	sessions, err := s.GetAllForUser(ctx, userID, nil)
	if err != nil {
		return errs.Internal(err)
	}

	for _, session := range sessions {
		err = s.Delete(ctx, &DeleteCmd{
			UserID: userID,
			Token:  session.Token(),
		})
		if err != nil {
			return errs.Internal(fmt.Errorf("failed to Delete web session %q: %w", session.Token(), err))
		}
	}

	return nil
}
