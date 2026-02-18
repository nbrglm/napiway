package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	LogoutUserRequestHTTPMethod = "GET"
	LogoutUserRequestRoutePath  = "/users/logout"
)

// Logout the current user.
type LogoutUserRequest struct {

	// Authentication parameters
	Auth LogoutUserRequestAuthParams
}

type LogoutUserRequestAuthParams struct {

	// Required Authentication Method
	// Source: header "X-App-API-Key"
	//
	// Authentication method that denotes an API key passed in the request header.
	//
	// Format (NOT ENFORCED): api_key
	//
	APIKey *string

	// At least one of the following authentication methods is required

	// Source: header "X-App-Session-Token"
	//
	// Authentication method that denotes a session token passed in the request header.
	//
	// Format (NOT ENFORCED): session_token
	//
	SessionToken *string

	// Source: header "X-App-Refresh-Token"
	//
	// Authentication method that denotes a refresh token passed in the request header.
	//
	// Format (NOT ENFORCED): refresh_token
	//
	RefreshToken *string
}

// NewLogoutUserRequest creates a new LogoutUserRequest from an http.Request and performs parameter parsing and validation.
func NewLogoutUserRequest(w http.ResponseWriter, r *http.Request) (req LogoutUserRequest, err error) {

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

	valRefreshToken := r.Header.Get("X-App-Refresh-Token")
	valRefreshToken = strings.TrimSpace(valRefreshToken)
	if valRefreshToken == "" {
		req.Auth.RefreshToken = nil
	} else {
		req.Auth.RefreshToken = &valRefreshToken
	}

	// Authentication parameters validation
	// Validate required auth parameters

	if req.Auth.APIKey == nil {
		err = fmt.Errorf("missing required authentication parameter: X-App-API-Key")
		return
	}

	authParamsSet := "X-App-Session-Token, X-App-Refresh-Token, "
	// Validate at least one auth parameter is present
	anyAuthParamsPresent := false

	if req.Auth.SessionToken != nil {
		anyAuthParamsPresent = true
	}

	if req.Auth.RefreshToken != nil {
		anyAuthParamsPresent = true
	}

	if !anyAuthParamsPresent {
		err = fmt.Errorf("missing required authentication (any of): %v", authParamsSet)
		return
	}

	return
}

type LogoutUser200Response struct {

	// Response body
	Body LogoutUser200ResponseBody
}

type LogoutUser200ResponseBody struct {

	// A message confirming successful logout.
	//
	// Required
	//
	// Must be non-empty
	Message string `json:"Message"`
}

// Successful logout response.
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteLogoutUser200Response(w http.ResponseWriter, response LogoutUser200Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(200)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type LogoutUser400Response struct {

	// Response body
	Body LogoutUser400ResponseBody
}

type LogoutUser400ResponseBodyError struct {

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

type LogoutUser400ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error LogoutUser400ResponseBodyError `json:"Error"`
}

// Bad Request
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteLogoutUser400Response(w http.ResponseWriter, response LogoutUser400Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(400)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type LogoutUser500Response struct {

	// Response body
	Body LogoutUser500ResponseBody
}

type LogoutUser500ResponseBodyError struct {

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

type LogoutUser500ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error LogoutUser500ResponseBodyError `json:"Error"`
}

// Internal Server Error
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteLogoutUser500Response(w http.ResponseWriter, response LogoutUser500Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(500)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}
