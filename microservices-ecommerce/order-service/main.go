package main

import (
	"encoding/json"
	"log"
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
	respondWithJSON(w, http.StatusOK, orders)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order, ok := orders[params["id"]]
	if !ok {
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}
	respondWithJSON(w, http.StatusOK, order)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if !isValidOrder(order) {
		respondWithError(w, http.StatusBadRequest, "Invalid order data")
		return
	}

	orders[order.ID] = order
	respondWithJSON(w, http.StatusCreated, order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedOrder Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if !isValidOrder(updatedOrder) {
		respondWithError(w, http.StatusBadRequest, "Invalid order data")
		return
	}

	orders[params["id"]] = updatedOrder
	respondWithJSON(w, http.StatusOK, updatedOrder)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := orders[params["id"]]; !ok {
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	delete(orders, params["id"])
	w.WriteHeader(http.StatusNoContent)
}

func isValidOrder(order Order) bool {
	return order.ID != "" && order.ProductID != "" && order.Quantity > 0 && order.Total > 0
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logAndRespondWithError(w, http.StatusInternalServerError, "Error encoding response")
	}
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	logAndRespondWithError(w, status, message)
}

func logAndRespondWithError(w http.ResponseWriter, status int, message string) {
	log.Printf("HTTP %d - %s", status, message)
	respondWithJSON(w, status, map[string]string{"error": message})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/orders", getOrders).Methods(http.MethodGet)
	r.HandleFunc("/orders/{id}", getOrder).Methods(http.MethodGet)
	r.HandleFunc("/orders", createOrder).Methods(http.MethodPost)
	r.HandleFunc("/orders/{id}", updateOrder).Methods(http.MethodPut)
	r.HandleFunc("/orders/{id}", deleteOrder).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8083", r))
}
