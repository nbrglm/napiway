package go_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type TestingAPI struct {
	baseURL string
	client  *http.Client
}

func NewTestingAPI(baseURL string) *TestingAPI {
	return &TestingAPI{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewTestingAPIWithClient(baseURL string, client *http.Client) *TestingAPI {
	return &TestingAPI{
		baseURL: baseURL,
		client:  client,
	}
}

func (c *TestingAPI) do(request *http.Request) (*http.Response, error) {
	if request.Header.Get("Accept") == "" {
		request.Header.Set("Accept", "application/json")
	}
	request.Header.Set("User-Agent", "TestingAPI-GoSDK/")
	return c.client.Do(request)
}

type CreateUserResult struct {
	StatusCode int

	Response201 CreateUser201Response

	Response400 CreateUser400Response

	Response500 CreateUser500Response

	UnknownResponse *UnknownStatusResponse
}

func (c *TestingAPI) CreateUser(ctx context.Context, params CreateUserRequest) (CreateUserResult, error) {
	if err := params.Validate(); err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid request parameters",
			Err:     err,
		}
	}

	body, err := json.Marshal(params.Body)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "failed to serialize request body",
			Err:     err,
		}
	}
	path := "/users/new"

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+path,
		bytes.NewReader(body),
	)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	authAPIKey, err := paramToString(params.Auth.APIKey, "auth parameter: APIKey", "*string", true)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid auth parameter: X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	authAdminToken, err := paramToString(params.Auth.AdminToken, "auth parameter: AdminToken", "*string", true)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid auth parameter: X-App-Admin-Token",
			Err:     err,
		}
	}
	req.Header.Set("X-App-Admin-Token", authAdminToken)

	resp, err := c.do(req)
	if err != nil {
		return CreateUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonNetworkError,
			Message: "network error during HTTP request",
			Err:     err,
		}
	}
	// resp.Body will be closed in response handlers
	response := CreateUserResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 201:
		result, err := NewCreateUser201Response(resp)
		if err != nil {
			return CreateUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response201 = result
		return response, nil

	case 400:
		result, err := NewCreateUser400Response(resp)
		if err != nil {
			return CreateUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response400 = result
		return response, nil

	case 500:
		result, err := NewCreateUser500Response(resp)
		if err != nil {
			return CreateUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response500 = result
		return response, nil

	default:
		response.UnknownResponse = &UnknownStatusResponse{
			StatusCode: resp.StatusCode,
			Response:   resp,
		}
		return response, nil
	}
}

type GetUserResult struct {
	StatusCode int

	Response200 GetUser200Response

	Response400 GetUser400Response

	Response404 GetUser404Response

	Response500 GetUser500Response

	UnknownResponse *UnknownStatusResponse
}

func (c *TestingAPI) GetUser(ctx context.Context, params GetUserRequest) (GetUserResult, error) {
	if err := params.Validate(); err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid request parameters",
			Err:     err,
		}
	}

	path := "/users/{userId}"

	pathParamUserId, err := paramToString(params.UserId, "path parameter: UserId", "string", true)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid path parameter: userId",
			Err:     err,
		}
	}
	path = strings.ReplaceAll(path, "{userId}", fmt.Sprintf("%v", pathParamUserId))

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+path,
		nil,
	)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	authAPIKey, err := paramToString(params.Auth.APIKey, "auth parameter: APIKey", "*string", true)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid auth parameter: X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	authSessionToken, err := paramToString(params.Auth.SessionToken, "auth parameter: SessionToken", "*string", true)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid auth parameter: X-App-Session-Token",
			Err:     err,
		}
	}
	req.Header.Set("X-App-Session-Token", authSessionToken)

	resp, err := c.do(req)
	if err != nil {
		return GetUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonNetworkError,
			Message: "network error during HTTP request",
			Err:     err,
		}
	}
	// resp.Body will be closed in response handlers
	response := GetUserResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:
		result, err := NewGetUser200Response(resp)
		if err != nil {
			return GetUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response200 = result
		return response, nil

	case 400:
		result, err := NewGetUser400Response(resp)
		if err != nil {
			return GetUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response400 = result
		return response, nil

	case 404:
		result, err := NewGetUser404Response(resp)
		if err != nil {
			return GetUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response404 = result
		return response, nil

	case 500:
		result, err := NewGetUser500Response(resp)
		if err != nil {
			return GetUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response500 = result
		return response, nil

	default:
		response.UnknownResponse = &UnknownStatusResponse{
			StatusCode: resp.StatusCode,
			Response:   resp,
		}
		return response, nil
	}
}

type HealthCheckResult struct {
	StatusCode int

	Response200 HealthCheck200Response

	UnknownResponse *UnknownStatusResponse
}

func (c *TestingAPI) HealthCheck(ctx context.Context, params HealthCheckRequest) (HealthCheckResult, error) {
	if err := params.Validate(); err != nil {
		return HealthCheckResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid request parameters",
			Err:     err,
		}
	}

	path := "/health"

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+path,
		nil,
	)
	if err != nil {
		return HealthCheckResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	resp, err := c.do(req)
	if err != nil {
		return HealthCheckResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonNetworkError,
			Message: "network error during HTTP request",
			Err:     err,
		}
	}
	// resp.Body will be closed in response handlers
	response := HealthCheckResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:
		result, err := NewHealthCheck200Response(resp)
		if err != nil {
			return HealthCheckResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response200 = result
		return response, nil

	default:
		response.UnknownResponse = &UnknownStatusResponse{
			StatusCode: resp.StatusCode,
			Response:   resp,
		}
		return response, nil
	}
}

type ListUsersResult struct {
	StatusCode int

	Response200 ListUsers200Response

	Response400 ListUsers400Response

	Response500 ListUsers500Response

	UnknownResponse *UnknownStatusResponse
}

func (c *TestingAPI) ListUsers(ctx context.Context, params ListUsersRequest) (ListUsersResult, error) {
	if err := params.Validate(); err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid request parameters",
			Err:     err,
		}
	}

	path := "/users"

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+path,
		nil,
	)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	authAPIKey, err := paramToString(params.Auth.APIKey, "auth parameter: APIKey", "*string", true)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid auth parameter: X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	authAdminToken, err := paramToString(params.Auth.AdminToken, "auth parameter: AdminToken", "*string", true)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid auth parameter: X-App-Admin-Token",
			Err:     err,
		}
	}
	req.Header.Set("X-App-Admin-Token", authAdminToken)

	q := req.URL.Query()

	queryPageNumber, err := paramToString(params.PageNumber, "query parameter: PageNumber", "*float64", false)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid query parameter: page",
			Err:     err,
		}
	}
	q.Set("page", queryPageNumber)

	queryPageSize, err := paramToString(params.PageSize, "query parameter: PageSize", "*float64", false)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid query parameter: pageSize",
			Err:     err,
		}
	}
	q.Set("pageSize", queryPageSize)

	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return ListUsersResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonNetworkError,
			Message: "network error during HTTP request",
			Err:     err,
		}
	}
	// resp.Body will be closed in response handlers
	response := ListUsersResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:
		result, err := NewListUsers200Response(resp)
		if err != nil {
			return ListUsersResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response200 = result
		return response, nil

	case 400:
		result, err := NewListUsers400Response(resp)
		if err != nil {
			return ListUsersResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response400 = result
		return response, nil

	case 500:
		result, err := NewListUsers500Response(resp)
		if err != nil {
			return ListUsersResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response500 = result
		return response, nil

	default:
		response.UnknownResponse = &UnknownStatusResponse{
			StatusCode: resp.StatusCode,
			Response:   resp,
		}
		return response, nil
	}
}

type LogoutUserResult struct {
	StatusCode int

	Response200 LogoutUser200Response

	Response400 LogoutUser400Response

	Response500 LogoutUser500Response

	UnknownResponse *UnknownStatusResponse
}

func (c *TestingAPI) LogoutUser(ctx context.Context, params LogoutUserRequest) (LogoutUserResult, error) {
	if err := params.Validate(); err != nil {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid request parameters",
			Err:     err,
		}
	}

	path := "/users/logout"

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+path,
		nil,
	)
	if err != nil {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	authAPIKey, err := paramToString(params.Auth.APIKey, "auth parameter: APIKey", "*string", true)
	if err != nil {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "invalid auth parameter: X-App-API-Key",
			Err:     err,
		}
	}
	req.Header.Set("X-App-API-Key", authAPIKey)

	numAuthParamsSet := 0

	authSessionToken, err := paramToString(params.Auth.SessionToken, "auth parameter: SessionToken", "*string", true)
	if err == nil {
		numAuthParamsSet++
		req.Header.Set("X-App-Session-Token", authSessionToken)
	}

	authRefreshToken, err := paramToString(params.Auth.RefreshToken, "auth parameter: RefreshToken", "*string", true)
	if err == nil {
		numAuthParamsSet++
		req.Header.Set("X-App-Refresh-Token", authRefreshToken)
	}

	if numAuthParamsSet != 1 {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonInvalidRequest,
			Message: "one auth parameter must be set",
			Err:     nil,
		}
	}

	resp, err := c.do(req)
	if err != nil {
		return LogoutUserResult{}, &TestingAPIError{
			Reason:  TestingAPIErrorReasonNetworkError,
			Message: "network error during HTTP request",
			Err:     err,
		}
	}
	// resp.Body will be closed in response handlers
	response := LogoutUserResult{
		StatusCode: resp.StatusCode,
	}
	switch resp.StatusCode {

	case 200:
		result, err := NewLogoutUser200Response(resp)
		if err != nil {
			return LogoutUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response200 = result
		return response, nil

	case 400:
		result, err := NewLogoutUser400Response(resp)
		if err != nil {
			return LogoutUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response400 = result
		return response, nil

	case 500:
		result, err := NewLogoutUser500Response(resp)
		if err != nil {
			return LogoutUserResult{}, &TestingAPIError{
				Reason:  TestingAPIErrorReasonDecodeError,
				Message: "failed to decode response",
				Err:     err,
			}
		}
		response.Response500 = result
		return response, nil

	default:
		response.UnknownResponse = &UnknownStatusResponse{
			StatusCode: resp.StatusCode,
			Response:   resp,
		}
		return response, nil
	}
}

type UnknownStatusResponse struct {
	StatusCode int
	Response   *http.Response
}
