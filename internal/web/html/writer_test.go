package html

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	t.Parallel()

	t.Run("WriteHTMLErrorPage success", func(t *testing.T) {
		html := NewRenderer(Config{
			PrettyRender: false,
			HotReload:    false,
		})

		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.RequestIDKey, "some-request-id")

		r := httptest.NewRequest(http.MethodGet, "/invalid-url", nil)
		r = r.WithContext(ctx)

		w := httptest.NewRecorder()

		html.WriteHTMLErrorPage(w, r, errors.New("some-error"))

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		rawBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		rawBodyStr := string(rawBody)

		assert.Contains(t, rawBodyStr, "RequestID: some-request-id")
		assert.Contains(t, rawBodyStr, "html lang=\"en\"", "must return the full page with the header")
	})
}
