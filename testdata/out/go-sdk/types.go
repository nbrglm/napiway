package go_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type TestingAPIErrorReason string

const (
	TestingAPIErrorReasonInvalidRequest TestingAPIErrorReason = "invalid_request"
	TestingAPIErrorReasonDecodeError    TestingAPIErrorReason = "decode_error"
	TestingAPIErrorReasonNetworkError   TestingAPIErrorReason = "network_error"
	TestingAPIErrorReasonOtherError     TestingAPIErrorReason = "other_error"
)

type TestingAPIError struct {
	Reason  TestingAPIErrorReason
	Message string
	Err     error
}

func (e *TestingAPIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Reason, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Reason, e.Message)
}

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

func (req *CreateUserRequest) Validate() error {

	if err := req.Body.Validate(); err != nil {
		return err
	}
	// Authentication parameters validation

	// Validate required auth parameters

	if req.Auth.APIKey == nil {
		return fmt.Errorf("missing required authentication parameter: X-App-API-Key")
	}
	if req.Auth.AdminToken == nil {
		return fmt.Errorf("missing required authentication parameter: X-App-Admin-Token")
	}

	return nil
}

func (o *CreateUserRequestBody) Validate() error {

	if strings.TrimSpace(o.Email) == "" {
		return fmt.Errorf("field 'Email' must be non-empty")
	}

	if strings.TrimSpace(o.UserName) == "" {
		return fmt.Errorf("field 'UserName' must be non-empty")
	}

	return nil
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

func NewCreateUser201Response(resp *http.Response) (CreateUser201Response, error) {
	defer resp.Body.Close()
	result := CreateUser201Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewCreateUser400Response(resp *http.Response) (CreateUser400Response, error) {
	defer resp.Body.Close()
	result := CreateUser400Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewCreateUser500Response(resp *http.Response) (CreateUser500Response, error) {
	defer resp.Body.Close()
	result := CreateUser500Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
}

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

func (req *GetUserRequest) Validate() error {

	// Authentication parameters validation

	// Validate required auth parameters

	if req.Auth.APIKey == nil {
		return fmt.Errorf("missing required authentication parameter: X-App-API-Key")
	}
	if req.Auth.SessionToken == nil {
		return fmt.Errorf("missing required authentication parameter: X-App-Session-Token")
	}

	return nil
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

func NewGetUser200Response(resp *http.Response) (GetUser200Response, error) {
	defer resp.Body.Close()
	result := GetUser200Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewGetUser400Response(resp *http.Response) (GetUser400Response, error) {
	defer resp.Body.Close()
	result := GetUser400Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewGetUser404Response(resp *http.Response) (GetUser404Response, error) {
	defer resp.Body.Close()
	result := GetUser404Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewGetUser500Response(resp *http.Response) (GetUser500Response, error) {
	defer resp.Body.Close()
	result := GetUser500Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
}

const (
	HealthCheckRequestHTTPMethod = "GET"
	HealthCheckRequestRoutePath  = "/health"
)

type HealthCheckRequest struct {
}

func (req *HealthCheckRequest) Validate() error {

	// Authentication parameters validation

	return nil
}

type HealthCheck200Response struct {

	// Response body
	Body HealthCheck200ResponseBody
}

type HealthCheck200ResponseBody struct {

	// Required
	//
	// Must be non-empty
	Status string `json:"Status"`
}

func NewHealthCheck200Response(resp *http.Response) (HealthCheck200Response, error) {
	defer resp.Body.Close()
	result := HealthCheck200Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
}

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

func (req *ListUsersRequest) Validate() error {

	// Authentication parameters validation

	// Validate required auth parameters

	if req.Auth.APIKey == nil {
		return fmt.Errorf("missing required authentication parameter: X-App-API-Key")
	}
	if req.Auth.AdminToken == nil {
		return fmt.Errorf("missing required authentication parameter: X-App-Admin-Token")
	}

	return nil
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

func NewListUsers200Response(resp *http.Response) (ListUsers200Response, error) {
	defer resp.Body.Close()
	result := ListUsers200Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewListUsers400Response(resp *http.Response) (ListUsers400Response, error) {
	defer resp.Body.Close()
	result := ListUsers400Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewListUsers500Response(resp *http.Response) (ListUsers500Response, error) {
	defer resp.Body.Close()
	result := ListUsers500Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
}

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

func (req *LogoutUserRequest) Validate() error {

	// Authentication parameters validation

	// Validate required auth parameters

	if req.Auth.APIKey == nil {
		return fmt.Errorf("missing required authentication parameter: X-App-API-Key")
	}

	// Get the set of auth parameters
	authParamsSet := "X-App-Session-Token, X-App-Refresh-Token"
	// Validate at least one auth parameter is present
	anyAuthParamsPresent := false

	if req.Auth.SessionToken != nil {
		anyAuthParamsPresent = true
	}

	if req.Auth.RefreshToken != nil {
		anyAuthParamsPresent = true
	}

	if !anyAuthParamsPresent {
		return fmt.Errorf("missing required authentication (any of): %v", authParamsSet)
	}
	return nil
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

func NewLogoutUser200Response(resp *http.Response) (LogoutUser200Response, error) {
	defer resp.Body.Close()
	result := LogoutUser200Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewLogoutUser400Response(resp *http.Response) (LogoutUser400Response, error) {
	defer resp.Body.Close()
	result := LogoutUser400Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
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

func NewLogoutUser500Response(resp *http.Response) (LogoutUser500Response, error) {
	defer resp.Body.Close()
	result := LogoutUser500Response{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result.Body); err != nil {
		return result, err
	}

	return result, nil
}
