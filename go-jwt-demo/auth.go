package main

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var privateKey *rsa.PrivateKey

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func loadPrivateKey() error {
	data, err := os.ReadFile("private.pem")
	if err != nil {
		return err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return err
	}

	privateKey = key
	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userID, role string

	// Fake authentication for now
	if req.Username == "admin" && req.Password == "password123" {
		userID = "123"
		role = "admin"
	} else if req.Username == "user" && req.Password == "password123" {
		userID = "456"
		role = "user"
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create jwt claims
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-demo",
			Subject:   userID,
			Audience:  []string{"go-jwt-api"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token using HS256 algo
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with our secret
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: signedToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
