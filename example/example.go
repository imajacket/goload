//go:build example

package example

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func InitDevRoutes(r *chi.Mux) {
	r.Get("/example", exampleHandler())
}

func exampleHandler() func(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("example/example.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Time string
		}{
			Time: time.Now().String(),
		}

		_ = tmpl.Execute(w, data)
	}
}
