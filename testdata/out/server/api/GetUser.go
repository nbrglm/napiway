package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// ParseGetUserReq creates a new instance of GetUserReq by parsing the http.Request
func ParseGetUserReq(w http.ResponseWriter, r *http.Request) (*GetUserReq, error) {
	req := GetUserReq{}
	var err error
	// to silence unused variable error in case there are no parameters to parse
	_ = err

	// Parse path parameters, if any

	var valUserId *string
	valUserId, err = parsestringParam(r.PathValue("userId"), "path: userId", true)
	if err != nil {
		return &GetUserReq{}, err
	}

	req.UserId = *valUserId

	// Parse query parameters, if any

	// Parse header parameters, if any

	// Required auth, if any

	valAPIKey := r.Header.Get("X-App-API-Key")
	valAPIKey = strings.TrimSpace(valAPIKey)
	if valAPIKey == "" {
		return &GetUserReq{}, fmt.Errorf("missing required authentication: header X-App-API-Key")
	} else {
		req.APIKeyAuth = valAPIKey
	}

	valSessionToken := r.Header.Get("X-App-Session-Token")
	valSessionToken = strings.TrimSpace(valSessionToken)
	if valSessionToken == "" {
		return &GetUserReq{}, fmt.Errorf("missing required authentication: header X-App-Session-Token")
	} else {
		req.SessionTokenAuth = valSessionToken
	}

	// Atleast one auth, if any

	return &req, nil
}

func NewGetUser200(

	body *User,

) *GetUser200 {
	return &GetUser200{

		Body: body,
	}
}

// Write200 writes the GetUser200 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *GetUserReq) Write200(w http.ResponseWriter, resp *GetUser200) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(200)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewGetUser400(

	body *ErrorResponse,

) *GetUser400 {
	return &GetUser400{

		Body: body,
	}
}

// Write400 writes the GetUser400 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *GetUserReq) Write400(w http.ResponseWriter, resp *GetUser400) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(400)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewGetUser404(

	body *ErrorResponse,

) *GetUser404 {
	return &GetUser404{

		Body: body,
	}
}

// Write404 writes the GetUser404 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *GetUserReq) Write404(w http.ResponseWriter, resp *GetUser404) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(404)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}

func NewGetUser500(

	body *ErrorResponse,

) *GetUser500 {
	return &GetUser500{

		Body: body,
	}
}

// Write500 writes the GetUser500 response to the http.ResponseWriter
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *GetUserReq) Write500(w http.ResponseWriter, resp *GetUser500) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(500)

	// Write body
	return json.NewEncoder(w).Encode(resp.Body)

}
