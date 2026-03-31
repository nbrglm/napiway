package go_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// NewLogoutUserReq creates a new instance of LogoutUserReq with required fields as parameters
func NewLogoutUserReq(

	APIKeyAuth string,

) *LogoutUserReq {
	return &LogoutUserReq{

		APIKeyAuth: APIKeyAuth,
	}
}

// WithRefreshTokenAuth sets the optional authentication parameter RefreshTokenAuth and returns the modified LogoutUserReq instance
func (o *LogoutUserReq) WithRefreshTokenAuth(value *string) *LogoutUserReq {
	o.RefreshTokenAuth = value
	return o
}

// WithSessionTokenAuth sets the optional authentication parameter SessionTokenAuth and returns the modified LogoutUserReq instance
func (o *LogoutUserReq) WithSessionTokenAuth(value *string) *LogoutUserReq {
	o.SessionTokenAuth = value
	return o
}

// ParseLogoutUser200 creates a new instance of LogoutUser200 by parsing a map[string]any
func ParseLogoutUser200(resp *http.Response) (*LogoutUser200, error) {
	result := new(LogoutUser200)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(LogoutUserResponseBody)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for LogoutUser200: %w", err)
	}

	return result, nil
}

// ParseLogoutUser400 creates a new instance of LogoutUser400 by parsing a map[string]any
func ParseLogoutUser400(resp *http.Response) (*LogoutUser400, error) {
	result := new(LogoutUser400)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for LogoutUser400: %w", err)
	}

	return result, nil
}

// ParseLogoutUser500 creates a new instance of LogoutUser500 by parsing a map[string]any
func ParseLogoutUser500(resp *http.Response) (*LogoutUser500, error) {
	result := new(LogoutUser500)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for LogoutUser500: %w", err)
	}

	return result, nil
}
