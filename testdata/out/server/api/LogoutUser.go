package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	LogoutUserReqHTTPMethod = "GET"
	LogoutUserReqRoutePath  = "/users/logout"
)

// Logout the current user.
type LogoutUserReq struct {

	// All of the below (upto AUTH-ALL-END comment) are required for authentication

	// Required Authentication Method
	// Source: header "X-App-API-Key"
	//
	// Authentication method that denotes an API key passed in the request header.
	//
	// Format (NOT ENFORCED): api_key
	//
	APIKeyAuth string

	// AUTH-ALL-END

	// At least one of the following (upto AUTH-ANY-END) authentication methods is required

	// Source: header "X-App-Refresh-Token"
	//
	// Authentication method that denotes a refresh token passed in the request header.
	//
	// Format (NOT ENFORCED): refresh_token
	//
	RefreshTokenAuth *string

	// Source: header "X-App-Session-Token"
	//
	// Authentication method that denotes a session token passed in the request header.
	//
	// Format (NOT ENFORCED): session_token
	//
	SessionTokenAuth *string

	// AUTH-ANY-END

	// NOTE: The RawBody field is not used here, as RequestBodyName and RawBody are mutually exclusive.
	// RawBody is only used in the golang client sdk generation, since that will make NewRequest functions more ergonomic to use for endpoints without a request body schema.
	// RawBody doesn't affect the structure of the request struct, and will not be unmarshalled/read by the server side NewRequest functions.
}

// Successful logout response.
type LogoutUser200 struct {

	// Response body
	Body *LogoutUserResponseBody
}

// Bad Request
type LogoutUser400 struct {

	// Response body
	Body *ErrorResponse
}

// Internal Server Error
type LogoutUser500 struct {

	// Response body
	Body *ErrorResponse
}

// ParseLogoutUserReq creates a new instance of LogoutUserReq by parsing the http.Request
func ParseLogoutUserReq(w http.ResponseWriter, r *http.Request) (*LogoutUserReq, error) {
	req := LogoutUserReq{}
	var err error
	// to silence unused variable error in case there are no parameters to parse
	_ = err

	// Parse path parameters, if any

	// Parse query parameters, if any

	// Parse header parameters, if any

	// Required auth, if any

	valAPIKey := r.Header.Get("X-App-API-Key")
	valAPIKey = strings.TrimSpace(valAPIKey)
	if valAPIKey == "" {
		return nil, fmt.Errorf("missing required authentication: header X-App-API-Key")
	} else {
		req.APIKeyAuth = valAPIKey
	}

	// Atleast one auth, if any

	anyAuthParamsPresent := false
	authParamsList := "X-App-Refresh-Token, X-App-Session-Token, "

	valRefreshToken := r.Header.Get("X-App-Refresh-Token")
	valRefreshToken = strings.TrimSpace(valRefreshToken)
	if valRefreshToken != "" {
		anyAuthParamsPresent = true
		req.RefreshTokenAuth = &valRefreshToken
	}

	valSessionToken := r.Header.Get("X-App-Session-Token")
	valSessionToken = strings.TrimSpace(valSessionToken)
	if valSessionToken != "" {
		anyAuthParamsPresent = true
		req.SessionTokenAuth = &valSessionToken
	}

	if !anyAuthParamsPresent {
		return nil, fmt.Errorf("at least one of the following authentication parameters is required: %s", authParamsList)
	}

	return &req, nil
}

func NewLogoutUser200(

	body *LogoutUserResponseBody,

) *LogoutUser200 {
	return &LogoutUser200{

		Body: body,
	}
}

// Write200 writes the LogoutUser200 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *LogoutUserReq) Write200(w http.ResponseWriter, resp *LogoutUser200) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(200)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewLogoutUser400(

	body *ErrorResponse,

) *LogoutUser400 {
	return &LogoutUser400{

		Body: body,
	}
}

// Write400 writes the LogoutUser400 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *LogoutUserReq) Write400(w http.ResponseWriter, resp *LogoutUser400) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(400)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewLogoutUser500(

	body *ErrorResponse,

) *LogoutUser500 {
	return &LogoutUser500{

		Body: body,
	}
}

// Write500 writes the LogoutUser500 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *LogoutUserReq) Write500(w http.ResponseWriter, resp *LogoutUser500) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(500)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}
