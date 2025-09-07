package internal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// DB wraps a Postgres connection pool
var DB *pgx.Conn

// ConnectDB initializes the database connection
func ConnectDB(connString string) error {
	var err error
	DB, err = pgx.Connect(context.Background(), connString)
	return err
}

// AuthenticateUser checks if the username and password are valid
func AuthenticateUser(username, password string) bool {
	var dbPassword string
	err := DB.QueryRow(context.Background(), "SELECT password FROM users WHERE username=$1", username).Scan(&dbPassword)
	if err != nil {
		return false
	}
	// Compare hashed password
	if bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)) != nil {
		return false
	}
	return true
}

// CreateUsersTable creates the users table if it does not exist
func CreateUsersTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		);
	`
	_, err := DB.Exec(context.Background(), query)
	return err
}

// CreateUser inserts a new user into the users table
func CreateUser(username, password string) error {
	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err = DB.Exec(context.Background(), query, username, string(hashed))
	return err
}
