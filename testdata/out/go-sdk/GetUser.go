package go_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	GetUserReqHTTPMethod = "GET"
	GetUserReqRoutePath  = "/users/{userId}"
)

// Retrieve user information by user ID.
type GetUserReq struct {

	// Source: path parameter "{userId}"
	//

	// The unique identifier of the user.
	//
	// Required
	UserId string

	// All of the below (upto AUTH-ALL-END comment) are required for authentication

	// Required Authentication Method
	// Source: header "X-App-API-Key"
	//
	// Authentication method that denotes an API key passed in the request header.
	//
	// Format (NOT ENFORCED): api_key
	//
	APIKeyAuth string

	// Required Authentication Method
	// Source: header "X-App-Session-Token"
	//
	// Authentication method that denotes a session token passed in the request header.
	//
	// Format (NOT ENFORCED): session_token
	//
	SessionTokenAuth string

	// AUTH-ALL-END

	// NOTE: The RawBody field is not used here, as RequestBodyName and RawBody are mutually exclusive.
	// RawBody is only used in the golang client sdk generation, since that will make NewRequest functions more ergonomic to use for endpoints without a request body schema.
	// RawBody doesn't affect the structure of the request struct, and will not be unmarshalled/read by the server side NewRequest functions.
}

// Successful response containing user information.
type GetUser200 struct {

	// Response body
	Body *User
}

// Bad Request
type GetUser400 struct {

	// Response body
	Body *ErrorResponse
}

// User Not Found
type GetUser404 struct {

	// Response body
	Body *ErrorResponse
}

// Internal Server Error
type GetUser500 struct {

	// Response body
	Body *ErrorResponse
}

// NewGetUserReq creates a new instance of GetUserReq with required fields as parameters
func NewGetUserReq(

	UserId string,

	APIKeyAuth string,

	SessionTokenAuth string,

) *GetUserReq {
	return &GetUserReq{

		UserId: UserId,

		APIKeyAuth: APIKeyAuth,

		SessionTokenAuth: SessionTokenAuth,
	}
}

// ParseGetUser200 creates a new instance of GetUser200 by parsing a map[string]any
func ParseGetUser200(resp *http.Response) (*GetUser200, error) {
	result := new(GetUser200)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(User)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for GetUser200: %w", err)
	}

	return result, nil
}

// ParseGetUser400 creates a new instance of GetUser400 by parsing a map[string]any
func ParseGetUser400(resp *http.Response) (*GetUser400, error) {
	result := new(GetUser400)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for GetUser400: %w", err)
	}

	return result, nil
}

// ParseGetUser404 creates a new instance of GetUser404 by parsing a map[string]any
func ParseGetUser404(resp *http.Response) (*GetUser404, error) {
	result := new(GetUser404)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for GetUser404: %w", err)
	}

	return result, nil
}

// ParseGetUser500 creates a new instance of GetUser500 by parsing a map[string]any
func ParseGetUser500(resp *http.Response) (*GetUser500, error) {
	result := new(GetUser500)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for GetUser500: %w", err)
	}

	return result, nil
}
