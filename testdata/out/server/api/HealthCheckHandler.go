package api

import (
	"encoding/json"
	"net/http"
)

const (
	HealthCheckRequestHTTPMethod = "GET"
	HealthCheckRequestRoutePath  = "/health"
)

type HealthCheckRequest struct {
}

// NewHealthCheckRequest creates a new HealthCheckRequest from an http.Request and performs parameter parsing and validation.
func NewHealthCheckRequest(w http.ResponseWriter, r *http.Request) (req HealthCheckRequest, err error) {

	return
}

type HealthCheck200Response struct {

	// Response body
	Body HealthCheck200ResponseBody
}

type HealthCheck200ResponseBody struct {

	// Required
	//
	// Must be non-empty
	Status string `json:"Status"`
}

// OK
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteHealthCheck200Response(w http.ResponseWriter, response HealthCheck200Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(200)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}
