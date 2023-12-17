package transport

import (
	"github.com/go-chi/render"
	"html/template"
	"log"
	"net/http"
)

func HomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("web/index.html")
		if err != nil {
			log.Printf("template parsing error: %v", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("template execute error: %v", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}
	}
}
