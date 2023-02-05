// templates
package main

import (
	"html/template"
	"path/filepath"

	"dandydev.com/todogo/pkg/models"
)

type templateData struct {
	Note     *models.Note
	Notes    []*models.Note
	Project  *models.Project
	Projects []*models.Project
	Status   *models.Status
	Statuses []*models.Status
}

func newTemplateCashe(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page)

		templates, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		templates, err = templates.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		templates, err = templates.ParseGlob(filepath.Join(dir, "*.partical.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[fileName] = templates
	}

	return cache, nil
}
