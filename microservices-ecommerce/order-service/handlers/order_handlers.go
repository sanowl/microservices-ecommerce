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

// NewOrderHandlers initializes the order handlers with the given orders map
func NewOrderHandlers() *OrderHandlers {
	orders := make(map[string]models.Order)
	orders["1"] = models.Order{ID: "1", ProductID: "101", Quantity: 1, Total: 100.0}
	orders["2"] = models.Order{ID: "2", ProductID: "102", Quantity: 2, Total: 200.0}
	return &OrderHandlers{Orders: orders}
}

// RegisterOrderHandlers registers the order handlers to the router
func (h *OrderHandlers) RegisterOrderHandlers(r *mux.Router) {
	r.HandleFunc("/orders", h.getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", h.getOrder).Methods("GET")
	r.HandleFunc("/orders", h.createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", h.updateOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}", h.deleteOrder).Methods("DELETE")
}

func (h *OrderHandlers) getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.Orders); err != nil {
		log.Printf("Error encoding orders: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *OrderHandlers) getOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order, ok := h.Orders[params["id"]]
	if ok {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			log.Printf("Error encoding order: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Order not found", http.StatusNotFound)
	}
}

func (h *OrderHandlers) createOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Printf("Error decoding order: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if order.ID == "" || order.ProductID == "" || order.Quantity <= 0 || order.Total <= 0 {
		http.Error(w, "Invalid order data", http.StatusBadRequest)
		return
	}
	h.Orders[order.ID] = order
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("Error encoding order: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *OrderHandlers) updateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		log.Printf("Error decoding order: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if _, ok := h.Orders[params["id"]]; !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	h.Orders[params["id"]] = updatedOrder
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedOrder); err != nil {
		log.Printf("Error encoding order: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *OrderHandlers) deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := h.Orders[params["id"]]; !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	delete(h.Orders, params["id"])
	w.WriteHeader(http.StatusNoContent)
}
