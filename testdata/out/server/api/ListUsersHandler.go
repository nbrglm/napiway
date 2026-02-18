package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	ListUsersRequestHTTPMethod = "GET"
	ListUsersRequestRoutePath  = "/users"
)

// List users with optional pagination.
type ListUsersRequest struct {

	// Source: query parameter "page"
	//

	// The page number for pagination. Default = 0.
	//
	// Optional
	PageNumber *float64

	// Source: query parameter "pageSize"
	//

	// The number of items per page for pagination. Default = 10.
	//
	// Optional
	PageSize *float64

	// Authentication parameters
	Auth ListUsersRequestAuthParams
}

type ListUsersRequestAuthParams struct {

	// Required Authentication Method
	// Source: header "X-App-API-Key"
	//
	// Authentication method that denotes an API key passed in the request header.
	//
	// Format (NOT ENFORCED): api_key
	//
	APIKey *string

	// Required Authentication Method
	// Source: header "X-App-Admin-Token"
	//
	// Authentication method that denotes an admin token passed in the request header.
	//
	// Format (NOT ENFORCED): admin_token
	//
	AdminToken *string
}

// NewListUsersRequest creates a new ListUsersRequest from an http.Request and performs parameter parsing and validation.
func NewListUsersRequest(w http.ResponseWriter, r *http.Request) (req ListUsersRequest, err error) {

	valPageNumber, err := parsefloat64Param(r.URL.Query().Get("page"), "query: page", false)
	if err != nil {
		return
	}

	req.PageNumber = valPageNumber

	valPageSize, err := parsefloat64Param(r.URL.Query().Get("pageSize"), "query: pageSize", false)
	if err != nil {
		return
	}

	req.PageSize = valPageSize

	valAPIKey := r.Header.Get("X-App-API-Key")
	valAPIKey = strings.TrimSpace(valAPIKey)
	if valAPIKey == "" {
		req.Auth.APIKey = nil
	} else {
		req.Auth.APIKey = &valAPIKey
	}

	valAdminToken := r.Header.Get("X-App-Admin-Token")
	valAdminToken = strings.TrimSpace(valAdminToken)
	if valAdminToken == "" {
		req.Auth.AdminToken = nil
	} else {
		req.Auth.AdminToken = &valAdminToken
	}

	// Authentication parameters validation
	// Validate required auth parameters

	if req.Auth.APIKey == nil {
		err = fmt.Errorf("missing required authentication parameter: X-App-API-Key")
		return
	}
	if req.Auth.AdminToken == nil {
		err = fmt.Errorf("missing required authentication parameter: X-App-Admin-Token")
		return
	}

	return
}

type ListUsers200Response struct {

	// Response body
	Body ListUsers200ResponseBody
}

type ListUsers200ResponseBodyUsersItem struct {

	// The age of the user.
	//
	// Optional
	//
	Age *float64 `json:"Age,omitempty"`

	// The email address of the user.
	//
	// Required
	//
	// Must be non-empty
	Email string `json:"Email"`

	// Indicates whether the user is active.
	//
	// Required
	//
	IsActive bool `json:"IsActive"`

	// The unique identifier of the user.
	//
	// Required
	//
	// Must be non-empty
	UserId string `json:"UserId"`

	// The name of the user.
	//
	// Required
	//
	// Must be non-empty
	UserName string `json:"UserName"`
}

type ListUsers200ResponseBody struct {

	// The current page number.
	//
	// Required
	//
	PageNumber float64 `json:"PageNumber"`

	// The number of items per page.
	//
	// Required
	//
	PageSize float64 `json:"PageSize"`

	// The total number of users available.
	//
	// Required
	//
	TotalCount float64 `json:"TotalCount"`

	// Required
	//
	Users []ListUsers200ResponseBodyUsersItem `json:"Users"`
}

// Successful response containing a list of users.
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteListUsers200Response(w http.ResponseWriter, response ListUsers200Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(200)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type ListUsers400Response struct {

	// Response body
	Body ListUsers400ResponseBody
}

type ListUsers400ResponseBodyError struct {

	// A detailed debug message for developers. Only passed if in debug mode.
	//
	// Optional
	//
	DebugMessage *string `json:"DebugMessage,omitempty"`

	// An error message which is user-friendly.
	//
	// Required
	//
	// Must be non-empty
	ErrorMessage string `json:"ErrorMessage"`
}

type ListUsers400ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error ListUsers400ResponseBodyError `json:"Error"`
}

// Bad Request
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteListUsers400Response(w http.ResponseWriter, response ListUsers400Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(400)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type ListUsers500Response struct {

	// Response body
	Body ListUsers500ResponseBody
}

type ListUsers500ResponseBodyError struct {

	// A detailed debug message for developers. Only passed if in debug mode.
	//
	// Optional
	//
	DebugMessage *string `json:"DebugMessage,omitempty"`

	// An error message which is user-friendly.
	//
	// Required
	//
	// Must be non-empty
	ErrorMessage string `json:"ErrorMessage"`
}

type ListUsers500ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error ListUsers500ResponseBodyError `json:"Error"`
}

// Internal Server Error
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteListUsers500Response(w http.ResponseWriter, response ListUsers500Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(500)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}
