package main

import (
	"html/template"
	"time"

	"github.com/lobre/doodle/pkg/embeds/htmldir"
	"github.com/lobre/doodle/pkg/forms"
	"github.com/lobre/doodle/pkg/models"

	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/shurcooL/httpfs/path/vfspath"
)

type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	Event           *models.Event
	Events          []*models.Event
}

// humanDate returns a nicely formatted string representation
// of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// functions holds custom functions that we want available
// in our templates.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache will load all template files, either from disk
// or from the embedded filesystem, and store them in an in-memory
// map for easy retrieval.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := vfspath.Glob(htmldir.FS, "*.page.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// inject our custom functions
		ts := template.New(page).Funcs(functions)

		ts, err = vfstemplate.ParseFiles(htmldir.FS, ts, page)
		if err != nil {
			return nil, err
		}

		ts, err = vfstemplate.ParseGlob(htmldir.FS, ts, "*.layout.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = vfstemplate.ParseGlob(htmldir.FS, ts, "*.partial.tmpl")
		if err != nil {
			return nil, err
		}

		cache[page] = ts
	}

	return cache, nil
}
