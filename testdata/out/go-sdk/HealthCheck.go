package go_sdk

import (
	"encoding/json"
	"fmt"
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

// NewHealthCheckReq creates a new instance of HealthCheckReq with required fields as parameters
func NewHealthCheckReq() *HealthCheckReq {
	return &HealthCheckReq{}
}

// ParseHealthCheck200 creates a new instance of HealthCheck200 by parsing a map[string]any
func ParseHealthCheck200(resp *http.Response) (*HealthCheck200, error) {
	result := new(HealthCheck200)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(HealthCheckResponseBody)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for HealthCheck200: %w", err)
	}

	return result, nil
}
