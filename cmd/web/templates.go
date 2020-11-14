package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/lobre/doodle/pkg/forms"
	"github.com/lobre/doodle/pkg/models"
)

type templateData struct {
	CurrentYear int
	Flash       string
	Form        *forms.Form
	Event       *models.Event
	Events      []*models.Event
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

// newTemplateCache will load all template files from disk into an in-memory
// map when the application starts. This speeds up the response time as we don't
// need to do it for all requests.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl').
		name := filepath.Base(page)

		// Inject our custom functions when creating the template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
