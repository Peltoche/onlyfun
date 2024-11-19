package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/Peltoche/onlyfun/internal/web/html/templates/auth"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Page_Bootstrap(t *testing.T) {
	t.Run("printPage success", func(t *testing.T) {
		htmlMock := html.NewMock(t)
		usersMock := users.NewMockService(t)
		handler := NewBootstrapPage(htmlMock, usersMock)

		// masterkeyMock.On("IsMasterKeyLoaded").Return(false).Once()

		htmlMock.On("WriteHTMLTemplate", mock.Anything, mock.Anything, http.StatusOK, &auth.BootstrapPageTmpl{}).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/bootstrap", nil)
		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("postForm success", func(t *testing.T) {
		htmlMock := html.NewMock(t)
		usersMock := users.NewMockService(t)
		handler := NewBootstrapPage(htmlMock, usersMock)

		newUser := users.NewFakeUser(t).Build()

		usersMock.On("Bootstrap", mock.Anything, &users.BootstrapCmd{
			Username: "username",
			Password: secret.NewText("some-secret"),
		}).Return(newUser, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/bootstrap", strings.NewReader(url.Values{
			"username": []string{"username"},
			"password": []string{"some-secret"},
			"confirm":  []string{"some-secret"},
		}.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusFound, res.StatusCode)
		assert.Equal(t, "/login", res.Header.Get("Location"))
	})

	t.Run("postForm with an invalid password confirmation", func(t *testing.T) {
		htmlMock := html.NewMock(t)
		usersMock := users.NewMockService(t)
		handler := NewBootstrapPage(htmlMock, usersMock)

		htmlMock.On("WriteHTMLTemplate", mock.Anything, mock.Anything, http.StatusOK, &auth.BootstrapPageTmpl{
			Username:      "username",
			Password:      secret.NewText("some-secret"),
			PasswordError: "",
			ConfirmError:  "not identical",
		}).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/bootstrap", strings.NewReader(url.Values{
			"username": []string{"username"},
			"password": []string{"some-secret"},
			"confirm":  []string{"not-the-same-secret"},
		}.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("postForm with a password too short", func(t *testing.T) {
		htmlMock := html.NewMock(t)
		usersMock := users.NewMockService(t)
		handler := NewBootstrapPage(htmlMock, usersMock)

		htmlMock.On("WriteHTMLTemplate", mock.Anything, mock.Anything, http.StatusOK, &auth.BootstrapPageTmpl{
			Username:      "username",
			Password:      secret.NewText("short"),
			PasswordError: "too short",
			ConfirmError:  "",
		}).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/bootstrap", strings.NewReader(url.Values{
			"username": []string{"username"},
			"password": []string{"short"},
			"confirm":  []string{"short"},
		}.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		srv := chi.NewRouter()
		handler.Register(srv, nil)
		srv.ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
