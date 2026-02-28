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

		baseContentTemplates := []string{"templates/base_content.gohtml"}
		if name == "base.gohtml" {
			continue
		} else if name == "login.gohtml" {
			baseContentTemplates = []string{}
		}

		patterns := append([]string{
			"templates/base.gohtml",
			"templates/disk_quota.gohtml",
			"templates/subcontent/**/*.gohtml",
			page,
		}, baseContentTemplates...)

		ts, err := template.New(name).Funcs(templateFunctions).ParseFS(internal.Templates, patterns...)

		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
