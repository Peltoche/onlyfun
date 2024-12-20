package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Templates(t *testing.T) {
	renderer := html.NewRenderer(html.Config{
		PrettyRender: false,
		HotReload:    false,
	})

	tests := []struct {
		Template html.Templater
		Name     string
		Layout   bool
	}{
		{
			Name:   "LoginPageTmpl",
			Layout: true,
			Template: &LoginPageTmpl{
				Username:      "some-user-input",
				UsernameError: "some-error-msg",
				PasswordError: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/foo", nil)

			if !test.Layout {
				r.Header.Add("HX-Boosted", "true")
			}

			renderer.WriteHTMLTemplate(w, r, http.StatusOK, test.Template)

			if !assert.Equal(t, http.StatusOK, w.Code) {
				res := w.Result()
				res.Body.Close()
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				t.Log(string(body))
			}
		})
	}
}
