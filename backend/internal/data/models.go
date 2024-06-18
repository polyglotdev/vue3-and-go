package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")

		}
		log.Printf("failed to get user by ID: %v", err)
		return nil, err
	}

	return &user, nil
}

// Update updates a user in the database.
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
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Printf("failed to generate hashed password: %v", err)
		return 0, err
	}

	var newID int
	query := "INSERT INTO users (email, first_name, last_name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) returning id"

	retryCount := 3
	for retries := 0; retries < retryCount; retries++ {
		err = db.QueryRowContext(ctx, query, user.Email, user.FirstName, user.LastName, hashedPassword, user.CreatedAt, user.UpdatedAt).Scan(&newID)
		if err == nil {
			user.ID = newID
			return 0, nil
		}
		log.Printf("failed to insert user, attempt %d: %v", retries+1, err)
		if retries < retryCount-1 {
			time.Sleep(time.Second * 2)
		}
	}
	return 0, fmt.Errorf("failed to insert user after %d attempts: %w", retryCount, err)
}

// ResetPassword resets the password for the user with the specified email address.
// It returns an error if the user does not exist or if there is an error resetting the password.
//
// Parameters:
//   - email: The email address of the user to reset the password for.
//
// Returns:
//   - An error if the user does not exist or if there is an error resetting the password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Select the user with the specified password address
	query := "update users set password = $1 where id = $2"
	_, err := db.ExecContext(ctx, query, password, u.ID)
	if err != nil {
		log.Printf("failed to update user: %v", err)
	}

	return nil
}

// PasswordMatches checks if the provided password matches the user's password.
// It returns true if the password matches, false otherwise.
//
// Parameters:
//   - password: The password to check.
//
// Returns:
//   - True if the password matches, false otherwise.
func (u *User) PasswordMatches(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
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

// GetByToken retrieves a token from the database by its token value.
// It returns a pointer to the Token struct if the token is found, or nil if not found.
//
// Parameters:
//   - token: The token value to retrieve.
//
// Returns:
//   - A pointer to the Token struct representing the token with the specified token value.
//   - An error if any occurs during the query execution or row scanning.
func (t *Token) GetByToken(plainText string) (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Select the plainText with the specified plainText value
	query := "SELECT id, user_id, email, plainText, token_hash, created_at, updated_at, expiry FROM tokens WHERE plainText = $1"

	var token Token
	row := db.QueryRowContext(ctx, query, plainText)
	err := row.Scan(&token.ID, &token.UserID, &token.Email, &token.Token, &token.TokenHash, &token.CreatedAt, &token.UpdatedAt, &token.Expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no plainText found with plainText %s", plainText)

		}
		log.Printf("failed to get plainText by plainText: %v", err)
		return nil, err
	}

	return &token, nil
}

// GetUserByToken retrieves the user associated with a token from the database.
// It returns a pointer to the User struct if the user is found, or nil if not found.
//
// Parameters:
//   - token: The token value to retrieve.
//
// Returns:
//   - A pointer to the User struct representing the user associated with the token.
//   - An error if any occurs during the query execution or row scanning.
func (t *Token) GetUserByToken(token Token) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Select the user associated with the token
	query := "SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE id = $1"

	var user User
	row := db.QueryRowContext(ctx, query, token.UserID)
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user found with token %v", token)

		}
		log.Printf("failed to get user by token (token user ID: %v): %v", token.UserID, err)
		return nil, err
	}

	return &user, nil
}

// GenerateToken generates a token for a user.
// It returns a pointer to the Token struct representing the generated token.
//
// Parameters:
//   - user: The User struct representing the user to generate a token for.
//
// Returns:
//   - A pointer to the Token struct representing the generated token.
//   - An error if any occurs during the token generation process.
func (t *Token) GenerateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Printf("failed to generate random bytes: %v", err)
		return nil, err
	}

	token.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Token))
	token.TokenHash = hash[:]

	return token, nil
}

// AuthenticationToken takes a pointer to http.Request and returns a pointer to User and an error.
// It returns a pointer to the User struct if the user is found, or nil if not found.
//
// Parameters:
//   - r: The http.Request struct representing the request.
//
// Returns:
//   - A pointer to the User struct representing the user with the specified email address.
//   - An error if any occurs during the query execution or row scanning.
func (t *Token) AuthenticationToken(r *http.Request) (*User, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("no token provided")
	}

	headerParts := strings.Split(token, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("invalid token format")
	}

	tk := headerParts[1]

	if len(tk) != 26 {
		return nil, errors.New("invalid token length")
	}

	tokenModel, err := t.GetByToken(tk)
	if err != nil {
		return nil, errors.New("no matching token found")
	}

	if tokenModel.Expiry.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	user, err := t.GetUserByToken(*tokenModel)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

// InsertToken inserts a new token into the database. It takes a token of
// type Token and returns an error if any occurs during the token insertion process.
//
// Parameters:
//   - token: The Token struct representing the token to insert.
//
// Returns:
//   - An error if any occurs during the token insertion process.
func (t *Token) InsertToken(token Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// delete the token if it exists
	err := t.DeleteToken(token.Token)
	if err != nil {
		log.Printf("failed to delete token: %v", err)
		return err
	}

	// insert the token
	query := "INSERT INTO tokens (user_id, email, token, token_hash, created_at, updated_at, expiry) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = db.ExecContext(ctx, query, token.UserID, token.Email, token.Token, token.TokenHash, time.Now(), time.Now(), token.Expiry)
	if err != nil {
		log.Printf("failed to insert token: %v", err)
		return err
	}

	return nil
}

// DeleteToken deletes a token from the database. It takes a token of type Token and returns an error if any occurs during the token deletion process.
//
// Parameters:
//   - token: The Token struct representing the token to delete.
//
// Returns:
//   - An error if any occurs during the token deletion process.
func (t *Token) DeleteToken(plainText string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "DELETE FROM tokens WHERE plainText = $1"
	_, err := db.ExecContext(ctx, query, plainText)
	if err != nil {
		log.Printf("failed to delete token: %v", err)
		return err
	}

	return nil
}

// GetUserWithToken retrieves the user associated with a token from the database.
// It returns a pointer to the User struct if the user is found, or nil if not found.
//
// Parameters:
//   - token: The token value to retrieve.
//
// Returns:
//   - A pointer to the User struct representing the user associated with the token.
//   - An error if any occurs during the query execution or row scanning.
func (t *Token) GetUserWithToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Select the user associated with the token
	query := "SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE id = $1"

	var user User
	row := db.QueryRowContext(ctx, query, token)
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user found with token %v", token)

		}
		log.Printf("failed to get user by token (token user ID: %v): %v", token, err)
		return nil, err
	}

	return &user, nil
}

// VaildateToken validates a token and returns a boolean indicating whether the token is valid or not.
//
// Parameters:
//   - token: The token value to validate.
//
// Returns:
//   - A boolean indicating whether the token is valid or not.
func (t *Token) VaildateToken(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Select the user associated with the token
	query := "SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE id = $1"

	var user User
	row := db.QueryRowContext(ctx, query, token)
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no user found with token %v", token)

		}
		log.Printf("failed to get user by token (token user ID: %v): %v", token, err)
		return false, err
	}

	return true, nil
}
