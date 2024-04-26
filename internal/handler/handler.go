package handler

import (
	"TestWebServer/internal/service"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// Handler struct holds the service instance used for handling requests.
type Handler struct {
	s *service.Service
}

// New creates a new instance of the Handler with the provided service.
func New(ser *service.Service) *Handler {
	return &Handler{s: ser}
}

// AskForOrder serves the HTML page for asking for an order.
func (h *Handler) AskForOrder(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./web/ask.html")
	t.Execute(w, nil)
}

// GetOrder handles the request to get an order by its ID.
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {

	// RequestBody struct defines the expected JSON structure for the request body.
	type RequestBody struct {
		ID string `json:"id"`
	}
	var requestBody RequestBody
	// Decode the request body into the RequestBody struct.
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("error unmarshaling ID", err)
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	log.Println("requested ", requestBody.ID)

	// Attempt to get the order using the provided ID.
	order, err := h.s.Get(requestBody.ID)
	if err != nil {
		log.Println("handler caught an error ", err)
		// Respond with a 404 error if the order is not found.
		RespondError(w, http.StatusNotFound, "id not found")
		return
	}
	log.Println("HANDLER LAYER: extracted from cash: ", string(order))
	// Respond with the order if found.
	RespondOrder(w, http.StatusOK, order)
}

// RespondOrder sends a JSON response with the order data.
func RespondOrder(w http.ResponseWriter, code int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

// RespondError sends a JSON response with an error message.
func RespondError(w http.ResponseWriter, code int, err string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
