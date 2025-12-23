package main

import (
	"emitrrconnect4/internal/server"
	"os"
	"log"
)

func main() {
	// Render environment variable se PORT uthata hai
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Local testing ke liye default port
	}

	log.Printf("Server starting on port %s...", port)

	// Aapka purana server start logic
	// Agar server.Start() block nahi karta hai, toh hume ListenAndServe handle karna hoga.
	// Lekin agar Start() function ke andar hi ListenAndServe hai, 
	// toh aapko internal/server folder ki file change karni padegi.
	
	server.Start() 
}
