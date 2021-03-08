package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

var templates *template.Template

func main() {
	r := chi.NewRouter()

	var err error
	templates, err = parseTemplates("templates/")
	if err != nil {
		log.Fatal(err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		executeTemplate(w, r, "home", nil)
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		executeTemplate(w, r, "about", nil)
	})

	r.Get("/contact", func(w http.ResponseWriter, r *http.Request) {
		executeTemplate(w, r, "contact", nil)
	})

	r.Post("/contact", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		executeTemplate(w, r, "contact_submit", map[string]interface{}{
			"Name": r.PostForm.Get("name"),
		})
	})

	http.ListenAndServe(":3000", r)
}

type page struct {
	Partial bool
	Data    interface{}
}

func executeTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	d := &page{
		Partial: r.Header.Get("HX-Request") == "true",
		Data:    data,
	}

	templates.ExecuteTemplate(w, name, d)
}

func parseTemplates(path string) (*template.Template, error) {
	tpl := template.New("")

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, ".html") {
			_, err = tpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tpl, nil
}
