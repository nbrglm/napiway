package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	GetUserRequestHTTPMethod = "GET"
	GetUserRequestRoutePath  = "/users/{userId}"
)

// Retrieve user information by user ID.
type GetUserRequest struct {

	// Source: path parameter "{userId}"
	//

	// The unique identifier of the user.
	//
	// Required
	UserId string

	// Authentication parameters
	Auth GetUserRequestAuthParams
}

type GetUserRequestAuthParams struct {

	// Required Authentication Method
	// Source: header "X-App-API-Key"
	//
	// Authentication method that denotes an API key passed in the request header.
	//
	// Format (NOT ENFORCED): api_key
	//
	APIKey *string

	// Required Authentication Method
	// Source: header "X-App-Session-Token"
	//
	// Authentication method that denotes a session token passed in the request header.
	//
	// Format (NOT ENFORCED): session_token
	//
	SessionToken *string
}

// NewGetUserRequest creates a new GetUserRequest from an http.Request and performs parameter parsing and validation.
func NewGetUserRequest(w http.ResponseWriter, r *http.Request) (req GetUserRequest, err error) {

	valUserId, err := parsestringParam(r.PathValue("userId"), "path: userId", true)
	if err != nil {
		return
	}

	req.UserId = *valUserId

	valAPIKey := r.Header.Get("X-App-API-Key")
	valAPIKey = strings.TrimSpace(valAPIKey)
	if valAPIKey == "" {
		req.Auth.APIKey = nil
	} else {
		req.Auth.APIKey = &valAPIKey
	}

	valSessionToken := r.Header.Get("X-App-Session-Token")
	valSessionToken = strings.TrimSpace(valSessionToken)
	if valSessionToken == "" {
		req.Auth.SessionToken = nil
	} else {
		req.Auth.SessionToken = &valSessionToken
	}

	// Authentication parameters validation
	// Validate required auth parameters

	if req.Auth.APIKey == nil {
		err = fmt.Errorf("missing required authentication parameter: X-App-API-Key")
		return
	}
	if req.Auth.SessionToken == nil {
		err = fmt.Errorf("missing required authentication parameter: X-App-Session-Token")
		return
	}

	return
}

type GetUser200Response struct {

	// Response body
	Body GetUser200ResponseBody
}

type GetUser200ResponseBodyUser struct {

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

type GetUser200ResponseBody struct {

	// Response Schema for GetUser endpoint.
	//
	// Required
	//
	User GetUser200ResponseBodyUser `json:"User"`
}

// Successful response containing user information.
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteGetUser200Response(w http.ResponseWriter, response GetUser200Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(200)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type GetUser400Response struct {

	// Response body
	Body GetUser400ResponseBody
}

type GetUser400ResponseBodyError struct {

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

type GetUser400ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error GetUser400ResponseBodyError `json:"Error"`
}

// Bad Request
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteGetUser400Response(w http.ResponseWriter, response GetUser400Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(400)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type GetUser404Response struct {

	// Response body
	Body GetUser404ResponseBody
}

type GetUser404ResponseBodyError struct {

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

type GetUser404ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error GetUser404ResponseBodyError `json:"Error"`
}

// User Not Found
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteGetUser404Response(w http.ResponseWriter, response GetUser404Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(404)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type GetUser500Response struct {

	// Response body
	Body GetUser500ResponseBody
}

type GetUser500ResponseBodyError struct {

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

type GetUser500ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error GetUser500ResponseBodyError `json:"Error"`
}

// Internal Server Error
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteGetUser500Response(w http.ResponseWriter, response GetUser500Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(500)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}
