package main

import (
	"fmt"
	"lacesGo/internal"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()
	// Connect to Postgres
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	if err := internal.ConnectDB(connString); err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	// Create users table if not exists
	if err := internal.CreateUsersTable(); err != nil {
		fmt.Println("Failed to create users table:", err)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Go microservice running!")
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, internal.HelloHandler())
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if data, err := internal.HealthCheckHandler(); err != nil {
			http.Error(w, "Health check failed", http.StatusInternalServerError)
		} else {
			w.Write(data)
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		if internal.AuthenticateUser(username, password) {
			secret := os.Getenv("JWT_SECRET")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": username,
				"exp":      time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, err := token.SignedString([]byte(secret))
			if err != nil {
				http.Error(w, "Could not generate token", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"token": "%s"}`, tokenString)
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			http.Error(w, "Username and password required", http.StatusBadRequest)
			return
		}
		err := internal.CreateUser(username, password)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "User registered successfully!")
	})
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
