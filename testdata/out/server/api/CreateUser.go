package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// ParseCreateUserReq creates a new instance of CreateUserReq by parsing the http.Request
func ParseCreateUserReq(w http.ResponseWriter, r *http.Request) (*CreateUserReq, error) {
	req := CreateUserReq{}
	var err error
	// to silence unused variable error in case there are no parameters to parse
	_ = err

	// Parse path parameters, if any

	// Parse query parameters, if any

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

	// Parse request body
	defer r.Body.Close()
	bodyData := make(map[string]any)

	maxBodyBytes := int64(256 << 10) // Default max body bytes: 256KB

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	err = json.NewDecoder(r.Body).Decode(&bodyData)
	if err != nil {
		return nil, fmt.Errorf("error parsing request body: %w", err)
	}
	var body *CreateUserRequestBody
	body, err = ParseCreateUserRequestBody(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error parsing request body: %w", err)
	}
	req.Body = body

	return &req, nil
}

func NewCreateUser201(

	body *CreateUserResponseBody,

) *CreateUser201 {
	return &CreateUser201{

		Body: body,
	}
}

// Write201 writes the CreateUser201 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *CreateUserReq) Write201(w http.ResponseWriter, resp *CreateUser201) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(201)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewCreateUser400(

	body *ErrorResponse,

) *CreateUser400 {
	return &CreateUser400{

		Body: body,
	}
}

// Write400 writes the CreateUser400 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *CreateUserReq) Write400(w http.ResponseWriter, resp *CreateUser400) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(400)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewCreateUser413Response() *CreateUser413Response {
	return &CreateUser413Response{}
}

// Write413 writes the CreateUser413Response response to the http.ResponseWriter
//
// RawBody is true, hence this function will only set the headers and write the status code, rest is to be done by the caller.
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *CreateUserReq) Write413(w http.ResponseWriter, resp *CreateUser413Response) error {
	// Set headers, if any

	// Set status code and write the header as there are no body to write
	w.WriteHeader(413)
	return nil

}

func NewCreateUser500(

	body *ErrorResponse,

) *CreateUser500 {
	return &CreateUser500{

		Body: body,
	}
}

// Write500 writes the CreateUser500 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *CreateUserReq) Write500(w http.ResponseWriter, resp *CreateUser500) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(500)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}
