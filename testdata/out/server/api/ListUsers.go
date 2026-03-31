package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// ParseListUsersReq creates a new instance of ListUsersReq by parsing the http.Request
func ParseListUsersReq(w http.ResponseWriter, r *http.Request) (*ListUsersReq, error) {
	req := ListUsersReq{}

	// Parse path parameters, if any

	// Parse query parameters, if any

	valPageNumber, err := parseint64Param(r.URL.Query().Get("page"), "query: page", false)
	if err != nil {
		return nil, err
	}

	req.PageNumber = valPageNumber

	valPageSize, err := parseint64Param(r.URL.Query().Get("pageSize"), "query: pageSize", false)
	if err != nil {
		return nil, err
	}

	req.PageSize = valPageSize

	// Parse header parameters, if any

	// Required auth, if any

	valAdminToken := r.Header.Get("X-App-Admin-Token")
	valAdminToken = strings.TrimSpace(valAdminToken)
	if valAdminToken == "" {
		return nil, fmt.Errorf("missing required authentication: header X-App-Admin-Token")
	} else {
		req.AdminTokenAuth = valAdminToken
	}

	valAPIKey := r.Header.Get("X-App-API-Key")
	valAPIKey = strings.TrimSpace(valAPIKey)
	if valAPIKey == "" {
		return nil, fmt.Errorf("missing required authentication: header X-App-API-Key")
	} else {
		req.APIKeyAuth = valAPIKey
	}

	// Atleast one auth, if any

	return &req, nil
}

func NewListUsers200(

	RateLimitRemaining int64,

	body *ListUsersResponseBody,

) *ListUsers200 {
	return &ListUsers200{

		RateLimitRemaining: RateLimitRemaining,

		Body: body,
	}
}

// Write200 writes the ListUsers200 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *ListUsersReq) Write200(w http.ResponseWriter, resp *ListUsers200) error {
	// Set headers, if any

	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%v", resp.RateLimitRemaining))

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(200)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewListUsers400(

	body *ErrorResponse,

) *ListUsers400 {
	return &ListUsers400{

		Body: body,
	}
}

// Write400 writes the ListUsers400 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *ListUsersReq) Write400(w http.ResponseWriter, resp *ListUsers400) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(400)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewListUsers500(

	body *ErrorResponse,

) *ListUsers500 {
	return &ListUsers500{

		Body: body,
	}
}

// Write500 writes the ListUsers500 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *ListUsersReq) Write500(w http.ResponseWriter, resp *ListUsers500) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(500)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}
