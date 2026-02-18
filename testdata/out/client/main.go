package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	sdk "github.com/nbrglm/napiway/testdata/out/go_sdk"
)

var (
	VALID_API_KEY         = "valid"
	VALID_SESSION_TOKEN   = "valid"
	VALID_REFRESH_TOKEN   = "valid"
	VALID_ADMIN_TOKEN     = "valid"
	INVALID_API_KEY       = "invalid"
	INVALID_SESSION_TOKEN = "invalid"
	INVALID_REFRESH_TOKEN = "invalid"
	INVALID_ADMIN_TOKEN   = "invalid"
	PAGE_SIZE             = 20.0
	AGE                   = 30.0
)

type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	IsActive bool     `json:"is_active"`
	Age      *float64 `json:"age"`
}

var age1 = 28.0

var users = []User{
	{
		ID:       "1",
		Name:     "Alice",
		Email:    "alice@example.com",
		IsActive: true,
		Age:      &age1,
	},
	{
		ID:       "2",
		Name:     "Bob",
		Email:    "bob@example.com",
		IsActive: false,
		Age:      nil,
	},
}

type Result map[string]bool

func main() {
	args := os.Args
	if len(args) < 2 || args[1] == "" {
		stdErr(true, "Please provide the server address as the first argument")
		return
	}
	serverAddr := args[1]
	api := sdk.NewTestingAPI(serverAddr)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result := Result{}

	// Test user logout
	logoutResult, err := testUserLogout(ctx, api)
	if err != nil {
		stdErr(false, "Test user logout failed: %v\n", err)
		return
	}

	// Store the result for printing later
	structToMapStringBool(logoutResult, &result, "LogoutUser")

	// Test list users
	listUsersResult, err := testListUsers(ctx, api)
	if err != nil {
		stdErr(false, "Test list users failed: %v\n", err)
		return
	}

	// Store the result for printing later
	structToMapStringBool(listUsersResult, &result, "ListUsers")

	// Test get user
	getUserResult, err := testGetUser(ctx, api)
	if err != nil {
		stdErr(false, "Test get user failed: %v\n", err)
		return
	}

	// Store the result for printing later
	structToMapStringBool(getUserResult, &result, "GetUser")

	// Test create user
	createUserResult, err := testCreateUser(ctx, api)
	if err != nil {
		stdErr(false, "Test create user failed: %v\n", err)
		return
	}

	// Store the result for printing later
	structToMapStringBool(createUserResult, &result, "CreateUser")

	// Print the final result
	printResult(result)
}

type UserLogoutResult struct {
	WithMissingAPIKey              bool
	WithMissingTokens              bool
	WithInvalidAPIKey              bool
	WithInvalidSessionToken        bool
	WithInvalidRefreshToken        bool
	ValidOperationWithSessionToken bool
	ValidOperationWithRefreshToken bool
}

func testUserLogout(ctx context.Context, api *sdk.TestingAPI) (result UserLogoutResult, err error) {
	_, err = api.LogoutUser(ctx, sdk.LogoutUserRequest{})
	if err != nil {
		if testingErr, ok := err.(*sdk.TestingAPIError); ok {
			if testingErr.Reason == sdk.TestingAPIErrorReasonInvalidRequest {
				result.WithMissingAPIKey = true
			} else {
				return result, err
			}
		} else {
			return result, err
		}
	}

	_, err = api.LogoutUser(ctx, sdk.LogoutUserRequest{
		Auth: sdk.LogoutUserRequestAuthParams{
			APIKey: &VALID_API_KEY,
		},
	})
	if err != nil {
		if testingErr, ok := err.(*sdk.TestingAPIError); ok {
			if testingErr.Reason == sdk.TestingAPIErrorReasonInvalidRequest {
				result.WithMissingTokens = true
			}
		} else {
			return result, err
		}
	}

	resInvalidAPIKey, err := api.LogoutUser(ctx, sdk.LogoutUserRequest{
		Auth: sdk.LogoutUserRequestAuthParams{
			APIKey:       &INVALID_API_KEY,
			SessionToken: &VALID_SESSION_TOKEN,
		},
	})
	if err != nil {
		return result, err
	}
	if resInvalidAPIKey.StatusCode == 400 {
		result.WithInvalidAPIKey = true
	}

	resInvalidRefreshToken, err := api.LogoutUser(ctx, sdk.LogoutUserRequest{
		Auth: sdk.LogoutUserRequestAuthParams{
			APIKey:       &VALID_API_KEY,
			RefreshToken: &INVALID_REFRESH_TOKEN,
		},
	})
	if err != nil {
		return result, err
	}
	if resInvalidRefreshToken.StatusCode == 400 {
		result.WithInvalidRefreshToken = true
	}

	resInvalidSessionToken, err := api.LogoutUser(ctx, sdk.LogoutUserRequest{
		Auth: sdk.LogoutUserRequestAuthParams{
			APIKey:       &VALID_API_KEY,
			SessionToken: &INVALID_SESSION_TOKEN,
		},
	})
	if err != nil {
		return result, err
	}
	if resInvalidSessionToken.StatusCode == 400 {
		result.WithInvalidSessionToken = true
	}

	resValidSessionToken, err := api.LogoutUser(ctx, sdk.LogoutUserRequest{
		Auth: sdk.LogoutUserRequestAuthParams{
			APIKey:       &VALID_API_KEY,
			SessionToken: &VALID_SESSION_TOKEN,
		},
	})
	if err != nil {
		return result, err
	}
	if resValidSessionToken.StatusCode == 200 && strings.Contains(resValidSessionToken.Response200.Body.Message, "SessionToken") {
		result.ValidOperationWithSessionToken = true
	}

	resValidRefreshToken, err := api.LogoutUser(ctx, sdk.LogoutUserRequest{
		Auth: sdk.LogoutUserRequestAuthParams{
			APIKey:       &VALID_API_KEY,
			RefreshToken: &VALID_REFRESH_TOKEN,
		},
	})
	if err != nil {
		stdErr(false, "LogoutUser request with valid refresh token failed: %v\n", err)
		return result, err
	}
	if resValidRefreshToken.StatusCode == 200 && strings.Contains(resValidRefreshToken.Response200.Body.Message, "RefreshToken") {
		result.ValidOperationWithRefreshToken = true
	}

	return result, nil
}

type ListUsersResult struct {
	WithMissingAPIKey                bool
	WithMissingAdminToken            bool
	WithInvalidAPIKey                bool
	WithInvalidAdminToken            bool
	ValidOperationWithoutQueryParams bool
	ValidOperationWithQueryParams    bool
}

func testListUsers(ctx context.Context, api *sdk.TestingAPI) (result ListUsersResult, err error) {
	_, err = api.ListUsers(ctx, sdk.ListUsersRequest{})
	if err != nil {
		if testingErr, ok := err.(*sdk.TestingAPIError); ok {
			if testingErr.Reason == sdk.TestingAPIErrorReasonInvalidRequest {
				result.WithMissingAPIKey = true
			} else {
				return result, err
			}
		} else {
			return result, err
		}
	}

	_, err = api.ListUsers(ctx, sdk.ListUsersRequest{
		Auth: sdk.ListUsersRequestAuthParams{
			APIKey: &VALID_API_KEY,
		},
	})
	if err != nil {
		if testingErr, ok := err.(*sdk.TestingAPIError); ok {
			if testingErr.Reason == sdk.TestingAPIErrorReasonInvalidRequest {
				result.WithMissingAdminToken = true
			}
		} else {
			return result, err
		}
	}

	resInvalidAPIKey, err := api.ListUsers(ctx, sdk.ListUsersRequest{
		Auth: sdk.ListUsersRequestAuthParams{
			APIKey:     &INVALID_API_KEY,
			AdminToken: &VALID_ADMIN_TOKEN,
		},
	})
	if err != nil {
		return result, err
	}
	if resInvalidAPIKey.StatusCode == 400 {
		result.WithInvalidAPIKey = true
	}

	resInvalidAdminToken, err := api.ListUsers(ctx, sdk.ListUsersRequest{
		Auth: sdk.ListUsersRequestAuthParams{
			APIKey:     &VALID_API_KEY,
			AdminToken: &INVALID_ADMIN_TOKEN,
		},
	})
	if err != nil {
		return result, err
	}
	if resInvalidAdminToken.StatusCode == 400 {
		result.WithInvalidAdminToken = true
	}

	resValidWithoutQueryParams, err := api.ListUsers(ctx, sdk.ListUsersRequest{
		Auth: sdk.ListUsersRequestAuthParams{
			APIKey:     &VALID_API_KEY,
			AdminToken: &VALID_ADMIN_TOKEN,
		},
	})
	if err != nil {
		return result, err
	}
	if resValidWithoutQueryParams.StatusCode == 200 {
		result.ValidOperationWithoutQueryParams = true
	}

	resValidWithQueryParams, err := api.ListUsers(ctx, sdk.ListUsersRequest{
		Auth: sdk.ListUsersRequestAuthParams{
			APIKey:     &VALID_API_KEY,
			AdminToken: &VALID_ADMIN_TOKEN,
		},
		PageSize: &PAGE_SIZE,
	})
	if err != nil {
		return result, err
	}
	if resValidWithQueryParams.StatusCode == 200 {
		result.ValidOperationWithQueryParams = true
	}
	return result, nil
}

type GetUserResult struct {
	WithoutPathParam bool
	ValidOperation   bool
}

func testGetUser(ctx context.Context, api *sdk.TestingAPI) (result GetUserResult, err error) {
	_, err = api.GetUser(ctx, sdk.GetUserRequest{
		Auth: sdk.GetUserRequestAuthParams{
			APIKey:       &VALID_API_KEY,
			SessionToken: &VALID_SESSION_TOKEN,
		},
	})
	if err != nil {
		if testingErr, ok := err.(*sdk.TestingAPIError); ok {
			if testingErr.Reason == sdk.TestingAPIErrorReasonInvalidRequest {
				result.WithoutPathParam = true
			} else {
				return result, err
			}
		} else {
			return result, err
		}
	}

	resValid, err := api.GetUser(ctx, sdk.GetUserRequest{
		Auth: sdk.GetUserRequestAuthParams{
			APIKey:       &VALID_API_KEY,
			SessionToken: &VALID_SESSION_TOKEN,
		},
		UserId: "1",
	})
	if err != nil {
		return result, err
	}
	if resValid.StatusCode == 200 && resValid.Response200.Body.User.UserId == "1" {
		result.ValidOperation = true
	}
	return result, nil
}

type CreateUserResult struct {
	ValidOperationWithOptionalFieldsMissing bool
	ValidOperationWithOptionalFieldPresent  bool
}

func testCreateUser(ctx context.Context, api *sdk.TestingAPI) (result CreateUserResult, err error) {
	resOptionalFieldMissing, err := api.CreateUser(ctx, sdk.CreateUserRequest{
		Auth: sdk.CreateUserRequestAuthParams{
			APIKey:     &VALID_API_KEY,
			AdminToken: &VALID_ADMIN_TOKEN,
		},
		Body: sdk.CreateUserRequestBody{
			Email:    "test@example.com",
			UserName: "Test User",
		},
	})
	if err != nil {
		return result, err
	}
	if resOptionalFieldMissing.StatusCode == 201 && resOptionalFieldMissing.Response201.Body.User.Email == "test@example.com" {
		result.ValidOperationWithOptionalFieldsMissing = true
	}

	resOptionalFieldPresent, err := api.CreateUser(ctx, sdk.CreateUserRequest{
		Auth: sdk.CreateUserRequestAuthParams{
			APIKey:     &VALID_API_KEY,
			AdminToken: &VALID_ADMIN_TOKEN,
		},
		Body: sdk.CreateUserRequestBody{
			Age:      &AGE,
			Email:    "test@example.com",
			UserName: "Test User",
		},
	})
	if err != nil {
		return result, err
	}
	if resOptionalFieldPresent.StatusCode == 201 && resOptionalFieldPresent.Response201.Body.User.Email == "test@example.com" {
		result.ValidOperationWithOptionalFieldPresent = true
	}

	return result, nil
}

func printResult(result Result) {
	marshalled, err := json.Marshal(result)
	if err != nil {
		stdErr(true, "Failed to marshal result: %v\n", err)
		return
	}
	os.Stdout.Write(marshalled)
}

func stdErr(exit bool, format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	if exit {
		os.Exit(1)
	}
}

func structToMapStringBool(input any, result *Result, prefix string) {
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	// Ensure we have a struct
	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Only process boolean fields
		if fieldValue.Kind() == reflect.Bool {
			(*result)[prefix+field.Name] = fieldValue.Bool()
		}
	}
}
