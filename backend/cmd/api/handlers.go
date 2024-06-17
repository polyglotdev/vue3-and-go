package main

import (
	"net/http"
)

// jsonResponse is the type used for generic JSON responses
// It is used to represent JSON responses that are sent back to the client.
// It has two fields:
//   - error: A boolean indicating whether an error occurred during the request.
//   - message: A string containing a message describing the result of the request.
type jsonResponse struct {
	// Error is a boolean indicating whether an error occurred during the request.
	Error bool `json:"error"`
	// Message is a string containing a message describing the result of the request.
	Message string `json:"message"`
}

// Login is the handler used to attempt to log a user into the api
// It expects a JSON object with the following fields:
//   - email: The email address of the user to log in.
//   - password: The password of the user to log in.
//
// It returns a JSON response with the following fields:
//   - error: A boolean indicating whether an error occurred during the login process.
//   - message: A string containing a message describing the result of the login process.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request.
//
// Returns:
//   - An error if there was an issue with the request or response.
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorLog.Println("Error while reading JSON:", err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	// TODO authenticate
	app.infoLog.Println(creds.Username, creds.Password)

	// send back a response
	payload.Error = false
	payload.Message = "Signed in"

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println("Error while writing JSON response:", err)
	}
}
