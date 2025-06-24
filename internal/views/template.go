package views

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

type Template struct {
	textTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.textTpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("failed to execute template %s: %v", name, err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// Create template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"seq": func(start, end int) []int {
			var result []int
			if start > end {
				for i := start; i >= end; i-- {
					result = append(result, i)
				}
			} else {
				for i := start; i <= end; i++ {
					result = append(result, i)
				}
			}
			return result
		},
		"slice": func(items interface{}, start, end int) interface{} {
			// This is a simplified slice function for templates
			return items
		},
	}

	tpl, err := template.New("").Funcs(funcMap).ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		textTpl: tpl,
	}, nil
}
