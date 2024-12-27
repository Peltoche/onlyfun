package assets

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func Test_Assets_HTTPHandler(t *testing.T) {
	t.Run("Handle success", func(t *testing.T) {
		handler := NewHTTPHandler(Config{HotReload: false, isTest: true})
		router := chi.NewRouter()
		handler.Register(router, nil)

		r := httptest.NewRequest(http.MethodGet, "/assets/hello.txt", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()

		// Status code
		require.Equal(t, http.StatusOK, res.StatusCode)

		// Body
		rawBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "Hello, World!\n", string(rawBody))

		// Headers
		require.NotEmpty(t, res.Header.Get("Content-Type"))
		require.NotEmpty(t, res.Header.Get("Last-Modified"))
		require.Equal(t, "14", res.Header.Get("Content-Length"))
	})

	t.Run("Handle with hot reload disable the cache", func(t *testing.T) {
		handler := NewHTTPHandler(Config{HotReload: true, isTest: true})
		router := chi.NewRouter()
		handler.Register(router, nil)

		r := httptest.NewRequest(http.MethodGet, "/assets/hello.txt", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()

		// Status code
		require.Equal(t, http.StatusOK, res.StatusCode)

		// Headers
		require.Equal(t, "no-cache", res.Header.Get("Cache-Control"))
	})
}
