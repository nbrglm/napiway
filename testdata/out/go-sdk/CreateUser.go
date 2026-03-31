package go_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CreateUserReqHTTPMethod = "POST"
	CreateUserReqRoutePath  = "/users/new"
)

// Create a new user in the system.
type CreateUserReq struct {

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

	// Request body
	Body *CreateUserRequestBody

	// NOTE: The RawBody field is not used here, as RequestBodyName and RawBody are mutually exclusive.
	// RawBody is only used in the golang client sdk generation, since that will make NewRequest functions more ergonomic to use for endpoints without a request body schema.
	// RawBody doesn't affect the structure of the request struct, and will not be unmarshalled/read by the server side NewRequest functions.
}

// Successful response containing the created user information.
type CreateUser201 struct {

	// Response body
	Body *CreateUserResponseBody
}

// Bad Request
type CreateUser400 struct {

	// Response body
	Body *ErrorResponse
}

// Payload Too Large - the request body exceeds the maximum allowed size
type CreateUser413Response struct {

	// Raw response body. The HTTP response will be returned directly for this response, and it will be the responsibility of the caller to read/close the response body.
	RawBody *http.Response
}

// Internal Server Error
type CreateUser500 struct {

	// Response body
	Body *ErrorResponse
}

// NewCreateUserReq creates a new instance of CreateUserReq with required fields as parameters
func NewCreateUserReq(

	AdminTokenAuth string,

	APIKeyAuth string,

	Body *CreateUserRequestBody,

) *CreateUserReq {
	return &CreateUserReq{

		AdminTokenAuth: AdminTokenAuth,

		APIKeyAuth: APIKeyAuth,

		Body: Body,
	}
}

// ParseCreateUser201 creates a new instance of CreateUser201 by parsing a map[string]any
func ParseCreateUser201(resp *http.Response) (*CreateUser201, error) {
	result := new(CreateUser201)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(CreateUserResponseBody)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for CreateUser201: %w", err)
	}

	return result, nil
}

// ParseCreateUser400 creates a new instance of CreateUser400 by parsing a map[string]any
func ParseCreateUser400(resp *http.Response) (*CreateUser400, error) {
	result := new(CreateUser400)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for CreateUser400: %w", err)
	}

	return result, nil
}

// ParseCreateUser413Response creates a new instance of CreateUser413Response by parsing a map[string]any
func ParseCreateUser413Response(resp *http.Response) (*CreateUser413Response, error) {
	result := new(CreateUser413Response)

	result.RawBody = resp

	return result, nil
}

// ParseCreateUser500 creates a new instance of CreateUser500 by parsing a map[string]any
func ParseCreateUser500(resp *http.Response) (*CreateUser500, error) {
	result := new(CreateUser500)

	defer resp.Body.Close()
	if result.Body == nil {
		result.Body = new(ErrorResponse)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result.Body); err != nil {
		return nil, fmt.Errorf("error decoding response body for CreateUser500: %w", err)
	}

	return result, nil
}
