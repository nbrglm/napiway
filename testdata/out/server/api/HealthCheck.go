package api

import (
	"encoding/json"
	"net/http"
)

const (
	HealthCheckReqHTTPMethod = "GET"
	HealthCheckReqRoutePath  = "/health"
)

type HealthCheckReq struct {

	// NOTE: The RawBody field is not used here, as RequestBodyName and RawBody are mutually exclusive.
	// RawBody is only used in the golang client sdk generation, since that will make NewRequest functions more ergonomic to use for endpoints without a request body schema.
	// RawBody doesn't affect the structure of the request struct, and will not be unmarshalled/read by the server side NewRequest functions.
}

// OK
type HealthCheck200 struct {

	// Response body
	Body *HealthCheckResponseBody
}

// ParseHealthCheckReq creates a new instance of HealthCheckReq by parsing the http.Request
func ParseHealthCheckReq(w http.ResponseWriter, r *http.Request) (*HealthCheckReq, error) {
	req := HealthCheckReq{}
	var err error
	// to silence unused variable error in case there are no parameters to parse
	_ = err

	// Parse path parameters, if any

	// Parse query parameters, if any

	// Parse header parameters, if any

	// Required auth, if any

	// Atleast one auth, if any

	return &req, nil
}

func NewHealthCheck200(

	body *HealthCheckResponseBody,

) *HealthCheck200 {
	return &HealthCheck200{

		Body: body,
	}
}

// Write200 writes the HealthCheck200 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *HealthCheckReq) Write200(w http.ResponseWriter, resp *HealthCheck200) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(200)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}
