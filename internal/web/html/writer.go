package html

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools/logger"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/render"
)

//go:embed *
var embeddedTemplates embed.FS

type Templater interface {
	Template() string
}

type Config struct {
	PrettyRender bool `mapstructure:"prettyRender"`
	HotReload    bool `mapstructure:"hotReload"`
}

type Writer interface {
	WriteHTMLTemplate(w http.ResponseWriter, r *http.Request, status int, template Templater)
	WriteHTMLErrorPage(w http.ResponseWriter, r *http.Request, err error)
}

type Renderer struct {
	render *render.Render
}

func NewRenderer(cfg Config) *Renderer {
	var directory string
	var fs render.FileSystem

	if cfg.HotReload {
		dir, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("failed to fetch the current workind dir: %s", err))
		}

		directory = path.Join(dir, "internal/web/html/templates")
		fs = render.LocalFileSystem{}
	} else {
		directory = ""
		fs = &render.EmbedFileSystem{FS: embeddedTemplates}
	}

	opts := render.Options{
		Directory:     directory,
		FileSystem:    fs,
		Layout:        "",
		IsDevelopment: cfg.HotReload,
		Extensions:    []string{".html"},
		Funcs: []template.FuncMap{
			{
				"humanTime": humanize.Time,
				"humanDate": func(t time.Time) string { return t.Format(time.DateTime) },
				"humanSize": humanize.Bytes,
			},
			{
				"sub": func(a, b int) int { return a - b },
			},
			{
				"add": func(a, b int) int { return a + b },
			},
			{
				"pathJoin": func(elems ...any) string {
					strElems := make([]string, len(elems))
					for i, elem := range elems {
						switch elem := elem.(type) {
						case uuid.UUID:
							strElems[i] = string(elem)
						default:
							strElems[i] = elem.(string)
						}
					}
					return path.Join(strElems...)
				},
				"getInodeIconClass": func(_ string, isDir bool) string {
					if isDir {
						return "fa-folder text-primary"
					}

					return "fa-file text-muted"
				},
			},
		},
	}

	if cfg.PrettyRender {
		opts.IndentXML = true
	}

	renderer := render.New(opts)
	renderer.CompileTemplates()

	return &Renderer{renderer}
}

func (t *Renderer) writeHTML(w http.ResponseWriter, r *http.Request, status int, template string, args any) {
	layout := ""

	if strings.Contains(template, "page") {
		dir := path.Dir(template)

		for {
			layout = path.Join(dir, "layout")
			if t.render.TemplateLookup(layout) != nil {
				break
			}

			dir = path.Dir(dir)

			if dir == "." {
				layout = ""
				break
			}
		}
	}

	if err := t.render.HTML(w, status, template, args, render.HTMLOptions{Layout: layout}); err != nil {
		logger.LogEntrySetAttrs(r.Context(), slog.String("render-error", err.Error()))
	}
}

func (t *Renderer) WriteHTMLTemplate(w http.ResponseWriter, r *http.Request, status int, template Templater) {
	t.writeHTML(w, r, status, template.Template(), template)
}

func (t *Renderer) WriteHTMLErrorPage(w http.ResponseWriter, r *http.Request, err error) {
	layout := ""

	reqID := r.Context().Value(middleware.RequestIDKey).(string)

	logger.LogEntrySetError(r.Context(), err)

	if err := t.render.HTML(w, http.StatusInternalServerError, "misc/page_500", map[string]any{
		"requestID": reqID,
	}, render.HTMLOptions{Layout: layout}); err != nil {
		logger.LogEntrySetAttrs(r.Context(), slog.String("render-error", err.Error()))
	}
}
