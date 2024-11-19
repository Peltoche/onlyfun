package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/services/websessions"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/auth"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_LoginPage(t *testing.T) {
	t.Parallel()

	t.Run("Login without any session open", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewMock(t)
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		handler := NewLoginPage(htmlMock, webSessionsMock, usersMock, tools)

		// Data

		// Mocks
		webSessionsMock.On("GetFromReq", mock.Anything).Return(nil, nil).Once()
		htmlMock.On("WriteHTMLTemplate", mock.Anything, mock.Anything, http.StatusOK, &auth.LoginPageTmpl{})

		// Run
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/login", nil)
		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		// Asserts
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("Login while already being authenticated redirect to the home", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewMock(t)
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		handler := NewLoginPage(htmlMock, webSessionsMock, usersMock, tools)

		// Data
		webSession := websessions.NewFakeSession(t).Build()

		// Mocks
		webSessionsMock.On("GetFromReq", mock.Anything).Return(webSession, nil).Once()

		// Run
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/login", nil)
		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		// Asserts
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusFound, res.StatusCode)
		assert.Equal(t, "/", res.Header.Get("Location"))
	})

	t.Run("ApplyLogin success", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewMock(t)
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		handler := NewLoginPage(htmlMock, webSessionsMock, usersMock, tools)

		// Data
		userPassword := gofakeit.Password(true, true, true, false, false, 8)
		user := users.NewFakeUser(t).WithPassword(userPassword).Build()
		webSession := websessions.NewFakeSession(t).
			CreatedBy(user).
			WithDevice("firefox 4.4.4.4").
			WithIP(httptest.DefaultRemoteAddr).
			Build()

		// Mocks
		usersMock.On("Authenticate", mock.Anything, user.Username(), secret.NewText(userPassword)).
			Return(user, nil).Once()
		webSessionsMock.On("Create", mock.Anything, &websessions.CreateCmd{
			UserID:     user.ID(),
			UserAgent:  "firefox 4.4.4.4",
			RemoteAddr: httptest.DefaultRemoteAddr,
		}).Return(webSession, nil).Once()

		// Run
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{
			"username": []string{user.Username()},
			"password": []string{userPassword},
		}.Encode()))
		r.RemoteAddr = httptest.DefaultRemoteAddr
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("User-Agent", "firefox 4.4.4.4")
		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		// Asserts
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusFound, res.StatusCode)
		assert.Equal(t, "/", res.Header.Get("Location"))
		assert.Len(t, res.Cookies(), 1)
		assert.Equal(t, "session_token", res.Cookies()[0].Name)
		assert.Empty(t, res.Cookies()[0].Expires)
	})

	t.Run("ApplyLogin success with the remember option", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewMock(t)
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		handler := NewLoginPage(htmlMock, webSessionsMock, usersMock, tools)

		// Data
		now := time.Now().UTC()
		userPassword := gofakeit.Password(true, true, true, false, false, 8)
		user := users.NewFakeUser(t).WithPassword(userPassword).Build()
		webSession := websessions.NewFakeSession(t).
			CreatedBy(user).
			WithDevice("firefox 4.4.4.4").
			WithIP(httptest.DefaultRemoteAddr).
			Build()

		// Mocks
		usersMock.On("Authenticate", mock.Anything, user.Username(), secret.NewText(userPassword)).
			Return(user, nil).Once()
		webSessionsMock.On("Create", mock.Anything, &websessions.CreateCmd{
			UserID:     user.ID(),
			UserAgent:  "firefox 4.4.4.4",
			RemoteAddr: httptest.DefaultRemoteAddr,
		}).Return(webSession, nil).Once()
		tools.ClockMock.On("Now").Return(now).Once()

		// Run
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{
			"username": []string{user.Username()},
			"password": []string{userPassword},
			"remember": []string{"on"},
		}.Encode()))
		r.RemoteAddr = httptest.DefaultRemoteAddr
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("User-Agent", "firefox 4.4.4.4")
		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		// Asserts
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusFound, res.StatusCode)
		assert.Equal(t, "/", res.Header.Get("Location"))
		assert.Len(t, res.Cookies(), 1)
		assert.Equal(t, "session_token", res.Cookies()[0].Name)
		assert.WithinDuration(t, now.Add(cookieLifeTime), res.Cookies()[0].Expires, 2*time.Second)
	})

	t.Run("ApplyLogin with an invalid username", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewMock(t)
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		handler := NewLoginPage(htmlMock, webSessionsMock, usersMock, tools)

		// Data

		// Mocks
		usersMock.On("Authenticate", mock.Anything, "invalid-username", secret.NewText("some-password")).
			Return(nil, users.ErrInvalidUsername).Once()
		htmlMock.On("WriteHTMLTemplate", mock.Anything, mock.Anything, http.StatusBadRequest, &auth.LoginPageTmpl{
			Username:      "invalid-username",
			UsernameError: "User doesn't exists",
			PasswordError: "",
		})

		// Run
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{
			"username": []string{"invalid-username"},
			"password": []string{"some-password"},
		}.Encode()))
		r.RemoteAddr = httptest.DefaultRemoteAddr
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("User-Agent", "firefox 4.4.4.4")

		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		// Asserts
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("ApplyLogin with an invalid password", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewMock(t)
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		handler := NewLoginPage(htmlMock, webSessionsMock, usersMock, tools)

		// Data
		user := users.NewFakeUser(t).Build()

		// Mocks
		usersMock.On("Authenticate", mock.Anything, user.Username(), secret.NewText("some-invalid-password")).
			Return(nil, users.ErrInvalidPassword).Once()
		htmlMock.On("WriteHTMLTemplate", mock.Anything, mock.Anything, http.StatusBadRequest, &auth.LoginPageTmpl{
			Username:      user.Username(),
			UsernameError: "",
			PasswordError: "Invalid password",
		})

		// Run
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{
			"username": []string{user.Username()},
			"password": []string{"some-invalid-password"},
		}.Encode()))
		r.RemoteAddr = httptest.DefaultRemoteAddr
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("User-Agent", "firefox 4.4.4.4")

		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		// Asserts
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("ApplyLogin with an authentication error", func(t *testing.T) {
		t.Parallel()

		tools := tools.NewMock(t)
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		handler := NewLoginPage(htmlMock, webSessionsMock, usersMock, tools)

		// Data
		user := users.NewFakeUser(t).Build()

		// Mocks
		usersMock.On("Authenticate", mock.Anything, user.Username(), secret.NewText("some-invalid-password")).
			Return(nil, errs.ErrInternal).Once()
		htmlMock.On("WriteHTMLErrorPage", mock.Anything, mock.Anything, errs.ErrInternal)

		// Run
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{
			"username": []string{user.Username()},
			"password": []string{"some-invalid-password"},
		}.Encode()))
		r.RemoteAddr = httptest.DefaultRemoteAddr
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("User-Agent", "firefox 4.4.4.4")

		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		// Asserts
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
