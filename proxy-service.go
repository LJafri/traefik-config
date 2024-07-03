package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func handleProxy(w http.ResponseWriter, r *http.Request) {
	targetService := r.Header.Get("x-service-endpoint")
	if targetService == "" {
		http.Error(w, "x-service-endpoint header is missing", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(targetService, "http://") && !strings.HasPrefix(targetService, "https://") {
		targetService = "http://" + targetService
	}

	targetURL := targetService + r.RequestURI
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}
	proxyReq.Header.Del("x-service-endpoint")

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

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
https://uat-us-api.experian.com/consumerservices/credit-profile/v1/decision-services
// curl -H "x-service-endpoint: uat-us-api.experian.com" "https://localhost:3030/consumerservices/credit-profile/v1/decision-services"                    

