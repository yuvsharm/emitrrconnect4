package server

import (
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("/ws", HandleWebSocket)

	log.Println("🚀 Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


