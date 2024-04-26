package handler

import "net/http"

// RegisterRoutes registers all the application routes.
func RegisterRoutes(h *Handler, mux *http.ServeMux) {
	mux.HandleFunc("/get_order", h.GetOrder)
	mux.HandleFunc("/", h.AskForOrder)
}
