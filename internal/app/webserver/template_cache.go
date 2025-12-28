package webserver

import (
	"html/template"
	"io/fs"
	"path/filepath"

	internal "github.com/almanac1631/scrubarr/web"
)

type TemplateCache map[string]*template.Template

func NewTemplateCache() (TemplateCache, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(internal.Templates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		if name == "base.html" {
			continue
		}

		ts, err := template.New(name).Funcs(templateFunctions).ParseFS(
			internal.Templates,
			"templates/base.gohtml",
			"templates/subcontent/**/*.gohtml",
			page,
		)

		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
