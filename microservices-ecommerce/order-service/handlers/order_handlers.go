package handlers

import (
	"encoding/json"
	"net/http"

	"order-service/models"

	"github.com/gorilla/mux"
)

var orders = make(map[string]models.Order)

func init() {
	// Initializing a few orders for demonstration
	orders["1"] = models.Order{ID: "1", ProductID: "101", Quantity: 1, Total: 100.0}
	orders["2"] = models.Order{ID: "2", ProductID: "102", Quantity: 2, Total: 200.0}
}

// RegisterOrderHandlers registers the order handlers to the router
func RegisterOrderHandlers(r *mux.Router) {
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", updateOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order, ok := orders[params["id"]]
	if ok {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Order not found", http.StatusNotFound)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orders[order.ID] = order
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, ok := orders[params["id"]]; !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	orders[params["id"]] = updatedOrder
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedOrder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := orders[params["id"]]; !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	delete(orders, params["id"])
	w.WriteHeader(http.StatusNoContent)
}
