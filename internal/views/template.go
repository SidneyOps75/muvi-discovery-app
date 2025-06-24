package views

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"reflect"
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
			// Use reflection to handle any slice type
			v := reflect.ValueOf(items)
			if v.Kind() != reflect.Slice {
				return items
			}
			
			length := v.Len()
			if start >= length {
				return reflect.MakeSlice(v.Type(), 0, 0).Interface()
			}
			if end > length {
				end = length
			}
			if start < 0 {
				start = 0
			}
			
			return v.Slice(start, end).Interface()
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
