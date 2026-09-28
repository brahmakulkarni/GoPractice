package main

import (
	"context"
	"crypto/rsa"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var publicKey *rsa.PublicKey

type contextKey string

const userClaimsKey contextKey = "userClaims"

func loadPublicKey() error {
	data, err := os.ReadFile("public.pem")
	if err != nil {
		return err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return err
	}

	publicKey = key
	return nil
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Expected format:
		// Authorization: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 3. Parse and validate JWT
		token, err := jwt.ParseWithClaims(
			tokenString,
			&Claims{},
			func(t *jwt.Token) (any, error) {
				// make sure the token uses the algo we expect
				if t.Method != jwt.SigningMethodRS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return publicKey, nil
			},
			jwt.WithIssuer("go-jwt-demo"),
			jwt.WithAudience("go-jwt-api"),
		)

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)

		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userClaimsKey).(*Claims)
		if !ok {
			http.Error(w, "User claims not found", http.StatusInternalServerError)
			return
		}

		if claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
