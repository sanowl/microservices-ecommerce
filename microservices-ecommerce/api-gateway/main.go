package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	target := vars["service"]

	var url string
	switch target {
	case "users":
		url = "http://user-service:8081" + r.URL.Path
	case "products":
		url = "http://product-service:8082" + r.URL.Path
	case "orders":
		url = "http://order-service:8083" + r.URL.Path
	default:
		http.NotFound(w, r)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{service}/{rest:.*}", proxyRequest)
	http.ListenAndServe(":8080", r)
}
