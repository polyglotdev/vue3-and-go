package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User:  User{},
		Token: Token{},
	}

}

type Models struct {
	User  User
	Token Token
}

// User is a struct that represents a user in the database
type User struct {
	// ID is the primary key of the user
	ID int `json:"id"`
	// Email is the email address of the user
	Email string `json:"email"`
	// FirstName is the first name of the user
	FirstName string `json:"first_name,omitempty"`
	// LastName is the last name of the user
	LastName string `json:"last_name,omitempty"`
	// Password is the password of the user
	Password string `json:"password"`
	// CreatedAt is the timestamp when the user was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the user was last updated
	UpdatedAt time.Time `json:"updated_at"`
	// Token is a struct that represents a token for the user
	Token Token `json:"token"`
}

// GetAll retrieves all users from the database.
// It returns a slice of User pointers and an error if any occurs during the query execution or row scanning.
//
// Returns:
//   - A slice of User pointers representing all users in the database.
//   - An error if any occurs during the query execution or row scanning.
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail retrieves a user from the database by their email address.
// It returns a pointer to the User struct if the user is found, or nil if not found.
//
// Parameters:
//   - email: The email address of the user to retrieve.
//
// Returns:
//   - A pointer to the User struct representing the user with the specified email address.
//   - An error if any occurs during the query execution or row scanning.
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Select the user with the specified email address
	query := "SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE email = $1"

	row := db.QueryRowContext(ctx, query, email)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email %s", email)

		}
		log.Printf("failed to get user by email: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user from the database by their ID.
// It returns a pointer to the User struct if the user is found, or nil if not found.
//
// Parameters:
//   - id: The ID of the user to retrieve.
//
// Returns:
//   - A pointer to the User struct representing the user with the specified ID.
//   - An error if any occurs during the query execution or row scanning.
func (u *User) GetByID(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Select the user with the specified ID
	query := "SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE id = $1"

	row := db.QueryRowContext(ctx, query, id)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")

		}
		log.Printf("failed to get user by ID: %v", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates a user in the database.
// It returns an error if any occurs during the query execution or row scanning.
//
// Parameters:
//   - user: The User struct representing the user to update.
//
// Returns:
//   - An error if any occurs during the query execution or row scanning.
func (u *User) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Update the user in the database
	query := "UPDATE users SET email = $1, first_name = $2, last_name = $3, updated_at = $4 WHERE id = $5"

	_, err := db.ExecContext(ctx, query, user.Email, user.FirstName, user.LastName, user.UpdatedAt, user.ID)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return err
	}

	return nil
}

// Delete deletes a user from the database.
// It returns an error if any occurs during the query execution or row scanning.
//
// Parameters:
//   - id: The ID of the user to delete.
//
// Returns:
//   - An error if any occurs during the query execution or row scanning.
func (u *User) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Delete the user from the database
	query := "DELETE FROM users WHERE id = $1"

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return err
	}

	return nil
}

// Insert inserts a user into the database.
// It returns an error if any occurs during the query execution or row scanning.
//
// Parameters:
//   - user: The User struct representing the user to insert.
//
// Returns:
//   - An error if any occurs during the query execution or row scanning.
func (u *User) Insert(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Insert the user into the database
	query := "INSERT INTO users (email, first_name, last_name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := db.ExecContext(ctx, query, user.Email, user.FirstName, user.LastName, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return err
	}

	return nil
}

// Token is a struct that represents a token in the database
type Token struct {
	// ID is the primary key of the token
	ID int `json:"id"`
	// UserID is the ID of the user associated with the token
	UserID int `json:"user_id"`
	// Email is the email address of the user associated with the token
	Email string `json:"email"`
	// Token is the token value
	Token string `json:"token"`
	// TokenHash is the hashed token value
	TokenHash []byte `json:"-"`
	// CreatedAt is the timestamp when the token was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the token was last updated
	UpdatedAt time.Time `json:"updated_at"`
	// Expiry is the expiry time of the token
	Expiry time.Time `json:"expiry"`
}
