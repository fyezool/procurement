package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/items", handleItems)
	http.HandleFunc("/quotes", handleQuotes)
	http.HandleFunc("/contracts", handleContracts)
	http.HandleFunc("/suppliers", handleSuppliers)
	http.HandleFunc("/procurements", handleProcurements)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering handleItems function")
	switch r.Method {
	case "GET":
		getItems(w, r)
	case "POST":
		createItem(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{ "message": "method not allowed" }`))
	}
	log.Println("Exiting handleItems function")
}

func handleQuotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getQuotes(w, r)
	case "POST":
		createQuote(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{ "message": "method not allowed" }`))
	}
}

func handleContracts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getContracts(w, r)
	case "POST":
		createContract(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{ "message": "method not allowed" }`))
	}
}

func handleSuppliers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getSuppliers(w, r)
	case "POST":
		createSupplier(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{ "message": "method not allowed" }`))
	}
}

func handleProcurements(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getProcurements(w, r)
	case "POST":
		createProcurement(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{ "message": "method not allowed" }`))
	}
}
