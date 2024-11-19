package language

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Peltoche/onlyfun/internal/tools/logger"
	"golang.org/x/text/language"
)

var browserLangCtxKey = &contextKey{"accept-language"}

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		browserLanguage := language.English

		browserLanguages, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		if err != nil {
			logger.LogEntrySetError(r.Context(), fmt.Errorf("failed to parse the accept-lanuage header: %w", err))
			// Set the default language
		}

		if len(browserLanguages) > 0 {
			browserLanguage = browserLanguages[0]
		}

		ctx := context.WithValue(r.Context(), browserLangCtxKey, browserLanguage)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// func SetBrowserLangFromReq(ctx context.Context, tag language.Tag) context.Context {
// 	return context.WithValue(ctx, browserLangCtxKey, tag)
// }
//
// func GetBrowserLangFromReq(ctx context.Context) language.Tag {
// 	return ctx.Value(browserLangCtxKey).(language.Tag)
// }

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "tools/language context value " + k.name
}
