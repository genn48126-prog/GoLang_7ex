package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	ledgerURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatalf("Invalid Ledger URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(ledgerURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Gateway received: %s %s", r.Method, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	log.Println("Gateway service starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}