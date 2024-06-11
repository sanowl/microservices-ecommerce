package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"order-service/models"

	"github.com/gorilla/mux"
)

// OrderHandlers holds the dependencies for the HTTP handlers
type OrderHandlers struct {
	Orders map[string]models.Order
}

// NewOrderHandlers initializes the order handlers with a predefined orders map
func NewOrderHandlers() *OrderHandlers {
	orders := map[string]models.Order{
		"1": {ID: "1", ProductID: "101", Quantity: 1, Total: 100.0},
		"2": {ID: "2", ProductID: "102", Quantity: 2, Total: 200.0},
	}
	return &OrderHandlers{Orders: orders}
}

// RegisterOrderHandlers registers the order handlers to the router
func (h *OrderHandlers) RegisterOrderHandlers(r *mux.Router) {
	r.HandleFunc("/orders", h.getOrders).Methods(http.MethodGet)
	r.HandleFunc("/orders/{id}", h.getOrder).Methods(http.MethodGet)
	r.HandleFunc("/orders", h.createOrder).Methods(http.MethodPost)
	r.HandleFunc("/orders/{id}", h.updateOrder).Methods(http.MethodPut)
	r.HandleFunc("/orders/{id}", h.deleteOrder).Methods(http.MethodDelete)
}

func (h *OrderHandlers) getOrders(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, h.Orders)
}

func (h *OrderHandlers) getOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order, ok := h.Orders[params["id"]]
	if !ok {
		h.respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}
	h.respondWithJSON(w, http.StatusOK, order)
}

func (h *OrderHandlers) createOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if !isValidOrder(order) {
		h.respondWithError(w, http.StatusBadRequest, "Invalid order data")
		return
	}

	h.Orders[order.ID] = order
	h.respondWithJSON(w, http.StatusCreated, order)
}

func (h *OrderHandlers) updateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if !isValidOrder(updatedOrder) {
		h.respondWithError(w, http.StatusBadRequest, "Invalid order data")
		return
	}

	existingOrder, ok := h.Orders[params["id"]]
	if !ok {
		h.respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	// Ensure the ID remains unchanged
	updatedOrder.ID = existingOrder.ID
	h.Orders[params["id"]] = updatedOrder
	h.respondWithJSON(w, http.StatusOK, updatedOrder)
}

func (h *OrderHandlers) deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := h.Orders[params["id"]]; !ok {
		h.respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	delete(h.Orders, params["id"])
	w.WriteHeader(http.StatusNoContent)
}

// isValidOrder validates the order data
func isValidOrder(order models.Order) bool {
	return order.ID != "" && order.ProductID != "" && order.Quantity > 0 && order.Total > 0
}

// respondWithJSON writes a JSON response to the client
func (h *OrderHandlers) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logAndRespondWithError(w, http.StatusInternalServerError, "Error encoding response")
	}
}

// respondWithError writes an error response to the client
func (h *OrderHandlers) respondWithError(w http.ResponseWriter, status int, message string) {
	h.logAndRespondWithError(w, status, message)
}

// logAndRespondWithError logs the error and writes an error response to the client
func (h *OrderHandlers) logAndRespondWithError(w http.ResponseWriter, status int, message string) {
	log.Printf("HTTP %d - %s", status, message)
	h.respondWithJSON(w, status, map[string]string{"error": message})
}
