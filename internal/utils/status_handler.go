package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// StatusHandler ...
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Server is up and running"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
}
