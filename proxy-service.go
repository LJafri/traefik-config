package main

import (
	"io"
	"log"
	"net/http"
)

func handleProxy(w http.ResponseWriter, r *http.Request) {
	targetService := r.Header.Get("x-service-endpoint")

	if targetService == "" {
		http.Error(w, "x-service-endpoint header is missing", http.StatusBadRequest)
		return
	}

	targetURL := "http://" + targetService + r.RequestURI

	// Create a new request
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}
	// Remove x-service-endpoint header for the upstream request
	proxyReq.Header.Del("x-service-endpoint")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", handleProxy)
	log.Println("Starting proxy service on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
