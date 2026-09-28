package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Go!")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)

	if !ok {
		http.Error(w, "User claims not found", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "User ID: %v\nRole: %v\n", claims.Subject, claims.Role)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the admin dashboard!")
}

func main() {
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	if err := loadPrivateKey(); err != nil {
		log.Fatal("Failed to load private key:", err)
	}

	if err := loadPublicKey(); err != nil {
		log.Fatal("Failed to load public key:", err)
	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/login", loginHandler)

	profile := http.HandlerFunc(profileHandler)
	http.Handle("/profile", authMiddleware(profile))

	admin := http.HandlerFunc(adminHandler)
	http.Handle("/admin", authMiddleware(
		adminMiddleware(admin),
	))

	fmt.Println("Server running on port http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
