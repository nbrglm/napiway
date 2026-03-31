package go_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const ClientVersion = "1.0.0"

type TestingAPIErrorReason string

const (
	// Network/Timeout
	ReasonTransport TestingAPIErrorReason = "transport"

	// Marshal/Unmarshal
	ReasonEncoding TestingAPIErrorReason = "encoding"

	// Non-Spec Response
	ReasonUnexpected TestingAPIErrorReason = "unexpected"
)

type TestingAPIError struct {
	// Reason is the "Actionable" category
	Reason TestingAPIErrorReason
	// Message is the human-readable "what happened"
	Message string
	// Err is the underlying cause (net.Conn, json.SyntaxError, etc.)
	Err error
}

func (e *TestingAPIError) Error() string {
	return fmt.Sprintf("%s_error: %s: %v", e.Reason, e.Message, e.Err)
}

type TestingAPI struct {
	httpClient *http.Client
	baseURL    string
}

func NewTestingAPI(baseURL string) *TestingAPI {
	return &TestingAPI{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
	}
}

func NewTestingAPIWithHTTPClient(baseURL string, httpClient *http.Client) *TestingAPI {
	return &TestingAPI{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *TestingAPI) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	req.Header.Set("User-Agent", "TestingAPI-GoSDK/1.0.0")
	return c.httpClient.Do(req.WithContext(ctx))
}

type CreateUserResult struct {

	// Successful response containing the created user information.
	Response201 *CreateUser201

	// Bad Request
	Response400 *CreateUser400

	// Payload Too Large - the request body exceeds the maximum allowed size
	Response413 *CreateUser413Response

	// Internal Server Error
	Response500 *CreateUser500

	StatusCode int

	// The raw http.Response for non-spec responses (e.g. 502 NGINX Error) that don't have a defined response body or headers in the spec. This allows callers to inspect the full response for debugging or error handling purposes.
	//
	// The client will only populate this field for responses that don't match any of the defined status codes in the spec.
	//
	// Callers should check the StatusCode field to determine if the response was a known spec response or an unknown response, and handle accordingly.
	//
	// Callers MUST CLOSE THE BODY of this response when done inspecting it to avoid resource leaks.
	UnknownResponse *http.Response
}

func (c *TestingAPI) CreateUser(ctx context.Context, params *CreateUserReq) (CreateUserResult, *TestingAPIError) {
	var body io.Reader

	bodyBytes, err := json.Marshal(params.Body)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "failed to marshal request body",
			Err:     err,
		}
	}
	body = bytes.NewReader(bodyBytes)

	path := "/users/new"

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+path,
		body,
	)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	authAdminToken, err := paramToString(params.AdminTokenAuth, "auth parameter: AdminToken", "string", true)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-Admin-Token",
			Err:     err,
		}
	}
	req.Header.Set("X-App-Admin-Token", authAdminToken)

	authAPIKey, err := paramToString(params.APIKeyAuth, "auth parameter: APIKey", "string", true)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	resp, err := c.do(ctx, req)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "HTTP request failed",
			Err:     err,
		}
	}
	response := CreateUserResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 201:

		parsedResp, err := ParseCreateUser201(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 201),
				Err:     err,
			}
		}
		response.Response201 = parsedResp
		return response, nil

	case 400:

		parsedResp, err := ParseCreateUser400(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 400),
				Err:     err,
			}
		}
		response.Response400 = parsedResp
		return response, nil

	case 413:

		parsedResp, err := ParseCreateUser413Response(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 413),
				Err:     err,
			}
		}
		response.Response413 = parsedResp
		return response, nil

	case 500:

		parsedResp, err := ParseCreateUser500(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 500),
				Err:     err,
			}
		}
		response.Response500 = parsedResp
		return response, nil

	default:
		response.UnknownResponse = resp
		return response, nil
	}
}

type GetUserResult struct {

	// Successful response containing user information.
	Response200 *GetUser200

	// Bad Request
	Response400 *GetUser400

	// User Not Found
	Response404 *GetUser404

	// Internal Server Error
	Response500 *GetUser500

	StatusCode int

	// The raw http.Response for non-spec responses (e.g. 502 NGINX Error) that don't have a defined response body or headers in the spec. This allows callers to inspect the full response for debugging or error handling purposes.
	//
	// The client will only populate this field for responses that don't match any of the defined status codes in the spec.
	//
	// Callers should check the StatusCode field to determine if the response was a known spec response or an unknown response, and handle accordingly.
	//
	// Callers MUST CLOSE THE BODY of this response when done inspecting it to avoid resource leaks.
	UnknownResponse *http.Response
}

func (c *TestingAPI) GetUser(ctx context.Context, params *GetUserReq) (GetUserResult, *TestingAPIError) {
	var body io.Reader

	path := "/users/{userId}"

	pathParamUserId, err := paramToString(params.UserId, "path parameter: UserId", "string", true)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid path parameter userId",
			Err:     err,
		}
	}
	path = strings.ReplaceAll(path, "{userId}", pathParamUserId)

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+path,
		body,
	)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	authAPIKey, err := paramToString(params.APIKeyAuth, "auth parameter: APIKey", "string", true)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	authSessionToken, err := paramToString(params.SessionTokenAuth, "auth parameter: SessionToken", "string", true)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-Session-Token",
			Err:     err,
		}
	}
	req.Header.Set("X-App-Session-Token", authSessionToken)

	resp, err := c.do(ctx, req)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "HTTP request failed",
			Err:     err,
		}
	}
	response := GetUserResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:

		parsedResp, err := ParseGetUser200(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 200),
				Err:     err,
			}
		}
		response.Response200 = parsedResp
		return response, nil

	case 400:

		parsedResp, err := ParseGetUser400(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 400),
				Err:     err,
			}
		}
		response.Response400 = parsedResp
		return response, nil

	case 404:

		parsedResp, err := ParseGetUser404(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 404),
				Err:     err,
			}
		}
		response.Response404 = parsedResp
		return response, nil

	case 500:

		parsedResp, err := ParseGetUser500(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 500),
				Err:     err,
			}
		}
		response.Response500 = parsedResp
		return response, nil

	default:
		response.UnknownResponse = resp
		return response, nil
	}
}

type ListUsersResult struct {

	// Successful response containing a list of users.
	Response200 *ListUsers200

	// Bad Request
	Response400 *ListUsers400

	// Internal Server Error
	Response500 *ListUsers500

	StatusCode int

	// The raw http.Response for non-spec responses (e.g. 502 NGINX Error) that don't have a defined response body or headers in the spec. This allows callers to inspect the full response for debugging or error handling purposes.
	//
	// The client will only populate this field for responses that don't match any of the defined status codes in the spec.
	//
	// Callers should check the StatusCode field to determine if the response was a known spec response or an unknown response, and handle accordingly.
	//
	// Callers MUST CLOSE THE BODY of this response when done inspecting it to avoid resource leaks.
	UnknownResponse *http.Response
}

func (c *TestingAPI) ListUsers(ctx context.Context, params *ListUsersReq) (ListUsersResult, *TestingAPIError) {
	var body io.Reader

	path := "/users"

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+path,
		body,
	)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	authAdminToken, err := paramToString(params.AdminTokenAuth, "auth parameter: AdminToken", "string", true)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-Admin-Token",
			Err:     err,
		}
	}
	req.Header.Set("X-App-Admin-Token", authAdminToken)

	authAPIKey, err := paramToString(params.APIKeyAuth, "auth parameter: APIKey", "string", true)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	q := req.URL.Query()

	queryPageNumber, err := paramToString(params.PageNumber, "query parameter: PageNumber", "*int64", false)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid query parameter page",
			Err:     err,
		}
	}
	q.Set("page", queryPageNumber)

	queryPageSize, err := paramToString(params.PageSize, "query parameter: PageSize", "*int64", false)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid query parameter pageSize",
			Err:     err,
		}
	}
	q.Set("pageSize", queryPageSize)

	req.URL.RawQuery = q.Encode()

	resp, err := c.do(ctx, req)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "HTTP request failed",
			Err:     err,
		}
	}
	response := ListUsersResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:

		parsedResp, err := ParseListUsers200(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 200),
				Err:     err,
			}
		}
		response.Response200 = parsedResp
		return response, nil

	case 400:

		parsedResp, err := ParseListUsers400(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 400),
				Err:     err,
			}
		}
		response.Response400 = parsedResp
		return response, nil

	case 500:

		parsedResp, err := ParseListUsers500(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 500),
				Err:     err,
			}
		}
		response.Response500 = parsedResp
		return response, nil

	default:
		response.UnknownResponse = resp
		return response, nil
	}
}

type LogoutUserResult struct {

	// Successful logout response.
	Response200 *LogoutUser200

	// Bad Request
	Response400 *LogoutUser400

	// Internal Server Error
	Response500 *LogoutUser500

	StatusCode int

	// The raw http.Response for non-spec responses (e.g. 502 NGINX Error) that don't have a defined response body or headers in the spec. This allows callers to inspect the full response for debugging or error handling purposes.
	//
	// The client will only populate this field for responses that don't match any of the defined status codes in the spec.
	//
	// Callers should check the StatusCode field to determine if the response was a known spec response or an unknown response, and handle accordingly.
	//
	// Callers MUST CLOSE THE BODY of this response when done inspecting it to avoid resource leaks.
	UnknownResponse *http.Response
}

func (c *TestingAPI) LogoutUser(ctx context.Context, params *LogoutUserReq) (LogoutUserResult, *TestingAPIError) {
	var body io.Reader

	path := "/users/logout"

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+path,
		body,
	)
	if err != nil {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	authAPIKey, err := paramToString(params.APIKeyAuth, "auth parameter: APIKey", "string", true)
	if err != nil {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	numAuthParamsSet := 0

	authRefreshToken, err := paramToString(params.RefreshTokenAuth, "auth parameter: RefreshToken", "*string", true)
	if err == nil {
		numAuthParamsSet++
		req.Header.Set("X-App-Refresh-Token", authRefreshToken)
	}

	authSessionToken, err := paramToString(params.SessionTokenAuth, "auth parameter: SessionToken", "*string", true)
	if err == nil {
		numAuthParamsSet++
		req.Header.Set("X-App-Session-Token", authSessionToken)
	}

	if numAuthParamsSet != 1 {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: fmt.Sprintf("exactly 1 auth parameter must be set, but %d were set", numAuthParamsSet),
			Err:     nil,
		}
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "HTTP request failed",
			Err:     err,
		}
	}
	response := LogoutUserResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:

		parsedResp, err := ParseLogoutUser200(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 200),
				Err:     err,
			}
		}
		response.Response200 = parsedResp
		return response, nil

	case 400:

		parsedResp, err := ParseLogoutUser400(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 400),
				Err:     err,
			}
		}
		response.Response400 = parsedResp
		return response, nil

	case 500:

		parsedResp, err := ParseLogoutUser500(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 500),
				Err:     err,
			}
		}
		response.Response500 = parsedResp
		return response, nil

	default:
		response.UnknownResponse = resp
		return response, nil
	}
}

// Invalid Request
//
// WhoAmI400 is a status-code only response.
type WhoAmIResult struct {

	// Successful response containing information about the currently authenticated user. Body is just a string with the user id provided in request body.
	Response200 *WhoAmI200

	StatusCode int

	// The raw http.Response for non-spec responses (e.g. 502 NGINX Error) that don't have a defined response body or headers in the spec. This allows callers to inspect the full response for debugging or error handling purposes.
	//
	// The client will only populate this field for responses that don't match any of the defined status codes in the spec.
	//
	// Callers should check the StatusCode field to determine if the response was a known spec response or an unknown response, and handle accordingly.
	//
	// Callers MUST CLOSE THE BODY of this response when done inspecting it to avoid resource leaks.
	UnknownResponse *http.Response
}

// NOTE: This endpoint has RawBody set to true, so the request body will not be handled by the generated client.
//
// Instead, the generated client function will have an additional parameter rawBody of type io.Reader, which will be the responsibility of the caller to read from and set the appropriate Content-Type header for the request.

func (c *TestingAPI) WhoAmI(ctx context.Context, params *WhoAmIReq, rawBody io.Reader) (WhoAmIResult, *TestingAPIError) {
	var body io.Reader

	body = rawBody

	path := "/users/whoami"

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+path,
		body,
	)
	if err != nil {
		return WhoAmIResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	authAPIKey, err := paramToString(params.APIKeyAuth, "auth parameter: APIKey", "string", true)
	if err != nil {
		return WhoAmIResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	authSessionToken, err := paramToString(params.SessionTokenAuth, "auth parameter: SessionToken", "string", true)
	if err != nil {
		return WhoAmIResult{}, &TestingAPIError{
			Reason:  ReasonEncoding,
			Message: "invalid auth parameter X-App-Session-Token",
			Err:     err,
		}
	}
	req.Header.Set("X-App-Session-Token", authSessionToken)

	resp, err := c.do(ctx, req)
	if err != nil {
		return WhoAmIResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "HTTP request failed",
			Err:     err,
		}
	}
	response := WhoAmIResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:

		parsedResp, err := ParseWhoAmI200(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 200),
				Err:     err,
			}
		}
		response.Response200 = parsedResp
		return response, nil

	case 400:

		// No response body or headers to parse for this status code
		return response, nil

	default:
		response.UnknownResponse = resp
		return response, nil
	}
}

type HealthCheckResult struct {

	// OK
	Response200 *HealthCheck200

	StatusCode int

	// The raw http.Response for non-spec responses (e.g. 502 NGINX Error) that don't have a defined response body or headers in the spec. This allows callers to inspect the full response for debugging or error handling purposes.
	//
	// The client will only populate this field for responses that don't match any of the defined status codes in the spec.
	//
	// Callers should check the StatusCode field to determine if the response was a known spec response or an unknown response, and handle accordingly.
	//
	// Callers MUST CLOSE THE BODY of this response when done inspecting it to avoid resource leaks.
	UnknownResponse *http.Response
}

func (c *TestingAPI) HealthCheck(ctx context.Context, params *HealthCheckReq) (HealthCheckResult, *TestingAPIError) {
	var body io.Reader

	path := "/health"

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+path,
		body,
	)
	if err != nil {
		return HealthCheckResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return HealthCheckResult{}, &TestingAPIError{
			Reason:  ReasonTransport,
			Message: "HTTP request failed",
			Err:     err,
		}
	}
	response := HealthCheckResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:

		parsedResp, err := ParseHealthCheck200(resp)
		if err != nil {
			response.UnknownResponse = resp
			return response, &TestingAPIError{
				Reason:  ReasonUnexpected,
				Message: fmt.Sprintf("failed to parse response for status code %d", 200),
				Err:     err,
			}
		}
		response.Response200 = parsedResp
		return response, nil

	default:
		response.UnknownResponse = resp
		return response, nil
	}
}
