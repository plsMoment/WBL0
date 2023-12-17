package transport

import (
	"WBL0/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"html/template"
	"log"
	"net/http"
)

func GetOrder(service service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")

		order, err := service.GetOrder(orderID)
		if err != nil {
			log.Printf("getting order failed: %v", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("web/order.html")
		if err != nil {
			log.Printf("parsing template failed: %v", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		if err = tmpl.Execute(w, string(order)); err != nil {
			log.Printf("executing template failed: %v", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}
	}
}
