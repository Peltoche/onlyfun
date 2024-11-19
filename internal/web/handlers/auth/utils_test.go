package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/services/websessions"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Utils_Authenticator(t *testing.T) {

	t.Run("getUserAndSession success", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		user := users.NewFakeUser(t).Build()
		session := websessions.NewFakeSession(t).CreatedBy(user).Build()

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(session, nil).Once()
		usersMock.On("GetByID", mock.Anything, user.ID()).Return(user, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		resUser, resSession, err := auth.GetUserAndSession(w, r)
		assert.Equal(t, user, resUser)
		assert.Equal(t, session, resSession)
		assert.Nil(t, err)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("getUserAndSession with a websession error", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("some-error")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		user, session, err := auth.GetUserAndSession(w, r)
		assert.Nil(t, user)
		assert.Nil(t, session)
		assert.ErrorIs(t, err, errs.ErrInternal)
	})

	t.Run("getUserAndSession with a websession not found", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(nil, websessions.ErrMissingSessionToken).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		user, session, err := auth.GetUserAndSession(w, r)
		assert.Nil(t, user)
		assert.Nil(t, session)
		assert.ErrorIs(t, err, ErrNotAuthenticated)
	})

	t.Run("getUserAndSession with a users problem", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		user := users.NewFakeUser(t).Build()
		session := websessions.NewFakeSession(t).CreatedBy(user).Build()

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(session, nil).Once()
		usersMock.On("GetByID", mock.Anything, user.ID()).Return(nil, fmt.Errorf("some-error")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		resUser, resSession, err := auth.GetUserAndSession(w, r)
		assert.Nil(t, resUser)
		assert.Nil(t, resSession)
		assert.ErrorIs(t, err, errs.ErrInternal)
	})

	t.Run("getUserAndSession with a user not found", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		user := users.NewFakeUser(t).Build()
		session := websessions.NewFakeSession(t).CreatedBy(user).Build()

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(session, nil).Once()
		usersMock.On("GetByID", mock.Anything, user.ID()).Return(nil, nil).Once()

		webSessionsMock.On("Logout", mock.Anything, mock.Anything).Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		user, session, err := auth.GetUserAndSession(w, r)
		assert.Nil(t, user)
		assert.Nil(t, session)
		assert.ErrorIs(t, err, ErrNotAuthenticated)
	})
}
