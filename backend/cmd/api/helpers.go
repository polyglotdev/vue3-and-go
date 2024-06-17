package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// readJSON reads and decodes JSON from an HTTP request body into the provided data structure.
// It ensures that the request body does not exceed a specified size limit and contains only a single JSON value.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request.
//   - data: A pointer to the data structure where the decoded JSON will be stored.
//
// Returns:
//   - An error if the JSON decoding fails or if the request body contains more than one JSON value.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

// writeJSON sends a JSON response with the specified status code and data.
// It also allows for optional HTTP headers to be included in the response.
//
// Parameters:
//
//	w       - The http.ResponseWriter to write the response to.
//	status  - The HTTP status code to set for the response.
//	data    - The data to be marshaled into JSON and included in the response body.
//	headers - Optional HTTP headers to include in the response.
//
// Returns:
//
//	error - An error if there was an issue marshaling the data or writing the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) errorJson(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	_ = app.writeJSON(w, statusCode, payload)
}
