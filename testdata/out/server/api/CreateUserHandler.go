package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	CreateUserRequestHTTPMethod = "POST"
	CreateUserRequestRoutePath  = "/users/new"
)

// Create a new user in the system.
type CreateUserRequest struct {

	// Authentication parameters
	Auth CreateUserRequestAuthParams

	// Request body
	Body CreateUserRequestBody
}

type CreateUserRequestBody struct {

	// The age of the user to be created.
	//
	// Optional
	//
	Age *float64 `json:"Age,omitempty"`

	// The email address of the user to be created.
	//
	// Required
	//
	// Must be non-empty
	Email string `json:"Email"`

	// The name of the user to be created.
	//
	// Required
	//
	// Must be non-empty
	UserName string `json:"UserName"`
}

type CreateUserRequestAuthParams struct {

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

func NewCreateUserRequestBody(data map[string]any) (CreateUserRequestBody, error) {
	var body CreateUserRequestBody

	valAge, ok := data["Age"]
	if !ok {

		// skip, leave as zero value

	} else {

		valAgeTyped, ok := valAge.(float64)
		if !ok {
			return body, fmt.Errorf("field 'Age' has incorrect type")
		}

		body.Age = &valAgeTyped

	}

	valEmail, ok := data["Email"]
	if !ok {

		return body, fmt.Errorf("missing required field 'Email'")

	} else {

		valEmailTyped, ok := valEmail.(string)
		if !ok {
			return body, fmt.Errorf("field 'Email' has incorrect type")
		}

		valEmailTyped = strings.TrimSpace(valEmailTyped)
		if len(valEmailTyped) == 0 {
			return body, fmt.Errorf("field 'Email' must be non-empty")
		}

		body.Email = valEmailTyped

	}

	valUserName, ok := data["UserName"]
	if !ok {

		return body, fmt.Errorf("missing required field 'UserName'")

	} else {

		valUserNameTyped, ok := valUserName.(string)
		if !ok {
			return body, fmt.Errorf("field 'UserName' has incorrect type")
		}

		valUserNameTyped = strings.TrimSpace(valUserNameTyped)
		if len(valUserNameTyped) == 0 {
			return body, fmt.Errorf("field 'UserName' must be non-empty")
		}

		body.UserName = valUserNameTyped

	}

	return body, nil
}

// NewCreateUserRequest creates a new CreateUserRequest from an http.Request and performs parameter parsing and validation.
func NewCreateUserRequest(w http.ResponseWriter, r *http.Request) (req CreateUserRequest, err error) {

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

	bodyData := make(map[string]any)

	maxBodyBytes := int64(256 << 10) // 256 KB default limit

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	err = json.NewDecoder(r.Body).Decode(&bodyData)
	if err != nil {
		return
	}
	var body CreateUserRequestBody
	body, err = NewCreateUserRequestBody(bodyData)
	if err != nil {
		return
	}
	req.Body = body
	defer r.Body.Close()

	return
}

type CreateUser201Response struct {

	// Response body
	Body CreateUser201ResponseBody
}

type CreateUser201ResponseBodyUser struct {

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

type CreateUser201ResponseBody struct {

	// Response Schema for GetUser endpoint.
	//
	// Required
	//
	User CreateUser201ResponseBodyUser `json:"User"`
}

// Successful response containing the created user information.
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteCreateUser201Response(w http.ResponseWriter, response CreateUser201Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(201)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type CreateUser400Response struct {

	// Response body
	Body CreateUser400ResponseBody
}

type CreateUser400ResponseBodyError struct {

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

type CreateUser400ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error CreateUser400ResponseBodyError `json:"Error"`
}

// Bad Request
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteCreateUser400Response(w http.ResponseWriter, response CreateUser400Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(400)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}

type CreateUser500Response struct {

	// Response body
	Body CreateUser500ResponseBody
}

type CreateUser500ResponseBodyError struct {

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

type CreateUser500ResponseBody struct {

	// Standard error response schema.
	//
	// Required
	//
	Error CreateUser500ResponseBodyError `json:"Error"`
}

// Internal Server Error
//
// This function WILL CALL w.WriteHeader(), so ensure that no other calls to
// w.WriteHeader() are made before calling this function.
func WriteCreateUser500Response(w http.ResponseWriter, response CreateUser500Response) error {
	// Set headers, if any

	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Set status code and write the header
	w.WriteHeader(500)

	// Write body
	return json.NewEncoder(w).Encode(response.Body)

}
