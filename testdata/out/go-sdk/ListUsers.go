package go_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ListUsersReqHTTPMethod = "GET"
	ListUsersReqRoutePath  = "/users"
)

// List users with optional pagination.
type ListUsersReq struct {

	// Source: query parameter "page"
	//

	// The page number for pagination. Default = 0.
	//
	// Optional
	PageNumber *int64

	// Source: query parameter "pageSize"
	//

	// The number of items per page for pagination. Default = 10.
	//
	// Optional
	PageSize *int64

	// All of the below (upto AUTH-ALL-END comment) are required for authentication

	// Required Authentication Method
	// Source: header "X-App-Admin-Token"
	//
	// Authentication method that denotes an admin token passed in the request header.
	//
	// Format (NOT ENFORCED): admin_token
	//
	AdminTokenAuth string

	// Required Authentication Method
	// Source: header "X-App-API-Key"
	//
	// Authentication method that denotes an API key passed in the request header.
	//
	// Format (NOT ENFORCED): api_key
	//
	APIKeyAuth string

	// AUTH-ALL-END

	// NOTE: The RawBody field is not used here, as RequestBodyName and RawBody are mutually exclusive.
	// RawBody is only used in the golang client sdk generation, since that will make NewRequest functions more ergonomic to use for endpoints without a request body schema.
	// RawBody doesn't affect the structure of the request struct, and will not be unmarshalled/read by the server side NewRequest functions.
}

// Successful response containing a list of users.
type ListUsers200 struct {

	// Source: header parameter "X-RateLimit-Remaining"
	//

	// The number of remaining requests allowed in the current rate limit window.
	//
	// Required
	RateLimitRemaining int64

	// Response body
	Body *ListUsersResponseBody
}

// Bad Request
type ListUsers400 struct {

	// Response body
	Body *ErrorResponse
}

// Internal Server Error
type ListUsers500 struct {

	// Response body
	Body *ErrorResponse
}

// NewListUsersReq creates a new instance of ListUsersReq with required fields as parameters
func NewListUsersReq(

	AdminTokenAuth string,

	APIKeyAuth string,

) *ListUsersReq {
	return &ListUsersReq{

		AdminTokenAuth: AdminTokenAuth,

		APIKeyAuth: APIKeyAuth,
	}
}

// WithPageNumber sets the optional query parameter PageNumber and returns the modified ListUsersReq instance
func (o *ListUsersReq) WithPageNumber(value *int64) *ListUsersReq {
	o.PageNumber = value
	return o
}

// WithPageSize sets the optional query parameter PageSize and returns the modified ListUsersReq instance
func (o *ListUsersReq) WithPageSize(value *int64) *ListUsersReq {
	o.PageSize = value
	return o
}

// ParseListUsers200 creates a new instance of ListUsers200 by parsing a map[string]any
func ParseListUsers200(resp *http.Response) (*ListUsers200, error) {
	result := new(ListUsers200)

	headerRateLimitRemaining, err := parseint64Param(resp.Header.Get("X-RateLimit-Remaining"), "header: X-RateLimit-Remaining", true)
	if err != nil {
		return nil, err
	}

	result.RateLimitRemaining = *headerRateLimitRemaining

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ListUsersResponseBody)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for ListUsers200: %w", err)
	}

	return result, nil
}

// ParseListUsers400 creates a new instance of ListUsers400 by parsing a map[string]any
func ParseListUsers400(resp *http.Response) (*ListUsers400, error) {
	result := new(ListUsers400)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for ListUsers400: %w", err)
	}

	return result, nil
}

// ParseListUsers500 creates a new instance of ListUsers500 by parsing a map[string]any
func ParseListUsers500(resp *http.Response) (*ListUsers500, error) {
	result := new(ListUsers500)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for ListUsers500: %w", err)
	}

	return result, nil
}
