package templates

import (
	"html/template"
	"net/http"
)

type TemplateInterface interface {
	ExecuteTemplate(w http.ResponseWriter, template string, data interface{}) error
}

type Template struct {
	Tmpl *template.Template
}

func (t *Template) LoadTemplates(pattern string) error {
	templates, err := template.ParseGlob(pattern)

	if err != nil {
		return err
	}

	t.Tmpl = templates
	return nil
}

func (t *Template) ExecuteTemplate(w http.ResponseWriter, template string, data interface{}) error {
	err := t.Tmpl.ExecuteTemplate(w, template, data)

	if err != nil {
		return err
	}

	return nil
}
