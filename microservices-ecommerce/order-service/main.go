package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Order struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

var orders = make(map[string]Order)

func getOrders(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(orders)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order, ok := orders[params["id"]]
	if ok {
		json.NewEncoder(w).Encode(order)
	} else {
		http.NotFound(w, r)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	_ = json.NewDecoder(r.Body).Decode(&order)
	orders[order.ID] = order
	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedOrder Order
	_ = json.NewDecoder(r.Body).Decode(&updatedOrder)
	orders[params["id"]] = updatedOrder
	json.NewEncoder(w).Encode(updatedOrder)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	delete(orders, params["id"])
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", updateOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")
	http.ListenAndServe(":8083", r)
}
