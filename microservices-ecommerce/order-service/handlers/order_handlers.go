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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.Orders); err != nil {
		h.handleError(w, http.StatusInternalServerError, "Error encoding orders", err)
	}
}

func (h *OrderHandlers) getOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order, ok := h.Orders[params["id"]]
	if !ok {
		h.handleError(w, http.StatusNotFound, "Order not found", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		h.handleError(w, http.StatusInternalServerError, "Error encoding order", err)
	}
}

func (h *OrderHandlers) createOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.handleError(w, http.StatusBadRequest, "Invalid order data", err)
		return
	}
	if !isValidOrder(order) {
		h.handleError(w, http.StatusBadRequest, "Invalid order fields", nil)
		return
	}

	h.Orders[order.ID] = order
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		h.handleError(w, http.StatusInternalServerError, "Error encoding order", err)
	}
}

func (h *OrderHandlers) updateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		h.handleError(w, http.StatusBadRequest, "Invalid order data", err)
		return
	}
	if !isValidOrder(updatedOrder) {
		h.handleError(w, http.StatusBadRequest, "Invalid order fields", nil)
		return
	}

	if _, ok := h.Orders[params["id"]]; !ok {
		h.handleError(w, http.StatusNotFound, "Order not found", nil)
		return
	}

	h.Orders[params["id"]] = updatedOrder
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedOrder); err != nil {
		h.handleError(w, http.StatusInternalServerError, "Error encoding updated order", err)
	}
}

func (h *OrderHandlers) deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := h.Orders[params["id"]]; !ok {
		h.handleError(w, http.StatusNotFound, "Order not found", nil)
		return
	}

	delete(h.Orders, params["id"])
	w.WriteHeader(http.StatusNoContent)
}

// isValidOrder validates the order data
func isValidOrder(order models.Order) bool {
	return order.ID != "" && order.ProductID != "" && order.Quantity > 0 && order.Total > 0
}

// handleError logs the error and sends an error response to the client
func (h *OrderHandlers) handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	} else {
		log.Println(message)
	}
	http.Error(w, message, statusCode)
}
