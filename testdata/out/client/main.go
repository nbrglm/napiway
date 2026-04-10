package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"reflect"
	"strings"
	"time"

	sdk "github.com/nbrglm/napiway/testdata/out/go_sdk"
)

var (
	VALID_API_KEY               = "valid"
	VALID_SESSION_TOKEN         = "valid"
	VALID_REFRESH_TOKEN         = "valid"
	VALID_ADMIN_TOKEN           = "valid"
	INVALID_API_KEY             = "invalid"
	INVALID_SESSION_TOKEN       = "invalid"
	INVALID_REFRESH_TOKEN       = "invalid"
	INVALID_ADMIN_TOKEN         = "invalid"
	PAGE_SIZE             int64 = 20
	AGE                   int64 = 30
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
	Age      *int64 `json:"age"`
}

var age1 int64 = 28

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

	// Test whoami
	whoAmIResult, err := testWhoAmI(ctx, api)
	if err != nil {
		stdErr(false, "Test whoami failed: %v\n", err)
		return
	}

	// Store the result for printing later
	structToMapStringBool(whoAmIResult, &result, "WhoAmI")

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

func testUserLogout(ctx context.Context, api *sdk.TestingAPI) (UserLogoutResult, error) {
	var result UserLogoutResult
	noApiKeyReq := sdk.NewLogoutUserReq("")
	_, err := api.LogoutUser(ctx, noApiKeyReq)
	if err != nil {
		if err.Reason == sdk.ReasonEncoding {
			result.WithMissingAPIKey = true
		} else {
			return result, err
		}
	}

	validApiKeyReq := sdk.NewLogoutUserReq(VALID_API_KEY)
	_, err = api.LogoutUser(ctx, validApiKeyReq)
	if err != nil {
		if err.Reason == sdk.ReasonEncoding {
			result.WithMissingTokens = true
		} else {
			return result, err
		}
	}

	invalidApiKeyReq := sdk.NewLogoutUserReq(INVALID_API_KEY).WithSessionTokenAuth(&VALID_SESSION_TOKEN)
	invalidApiKeyResp, err := api.LogoutUser(ctx, invalidApiKeyReq)
	if err != nil {
		return result, err
	}
	if invalidApiKeyResp.StatusCode == 400 {
		result.WithInvalidAPIKey = true
	}

	invalidRefreshTokenReq := sdk.NewLogoutUserReq(VALID_API_KEY).WithRefreshTokenAuth(&INVALID_REFRESH_TOKEN)
	invalidRefreshTokenResp, err := api.LogoutUser(ctx, invalidRefreshTokenReq)
	if err != nil {
		return result, err
	}
	if invalidRefreshTokenResp.StatusCode == 400 {
		result.WithInvalidRefreshToken = true
	}

	invalidSessionTokenReq := sdk.NewLogoutUserReq(VALID_API_KEY).WithSessionTokenAuth(&INVALID_SESSION_TOKEN)
	invalidSessionTokenResp, err := api.LogoutUser(ctx, invalidSessionTokenReq)
	if err != nil {
		return result, err
	}
	if invalidSessionTokenResp.StatusCode == 400 {
		result.WithInvalidSessionToken = true
	}

	validSessionTokenReq := sdk.NewLogoutUserReq(VALID_API_KEY).WithSessionTokenAuth(&VALID_SESSION_TOKEN)
	validSessionTokenResp, err := api.LogoutUser(ctx, validSessionTokenReq)
	if err != nil {
		return result, err
	}
	if validSessionTokenResp.StatusCode == 200 && strings.Contains(validSessionTokenResp.Response200.Body.Message, "SessionToken") {
		result.ValidOperationWithSessionToken = true
	}

	validRefreshTokenReq := sdk.NewLogoutUserReq(VALID_API_KEY).WithRefreshTokenAuth(&VALID_REFRESH_TOKEN)
	validRefreshTokenResp, err := api.LogoutUser(ctx, validRefreshTokenReq)
	if err != nil {
		return result, err
	}
	if validRefreshTokenResp.StatusCode == 200 && strings.Contains(validRefreshTokenResp.Response200.Body.Message, "RefreshToken") {
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

func testListUsers(ctx context.Context, api *sdk.TestingAPI) (ListUsersResult, error) {
	var result ListUsersResult

	invalidReq := sdk.NewListUsersReq("", "")
	_, err := api.ListUsers(ctx, invalidReq)
	if err != nil {
		if err.Reason == sdk.ReasonEncoding {
			result.WithMissingAPIKey = true
		} else {
			return result, err
		}
	}

	validApiKeyReq := sdk.NewListUsersReq("", VALID_API_KEY)
	_, err = api.ListUsers(ctx, validApiKeyReq)
	if err != nil {
		if err.Reason == sdk.ReasonEncoding {
			result.WithMissingAdminToken = true
		} else {
			return result, err
		}
	}

	reqInvalidAPIKey := sdk.NewListUsersReq(VALID_ADMIN_TOKEN, INVALID_API_KEY)
	resInvalidAPIKey, err := api.ListUsers(ctx, reqInvalidAPIKey)
	if err != nil {
		return result, err
	}
	if resInvalidAPIKey.StatusCode == 400 {
		result.WithInvalidAPIKey = true
	}

	reqInvalidAdminToken := sdk.NewListUsersReq(INVALID_ADMIN_TOKEN, VALID_API_KEY)
	resInvalidAdminToken, err := api.ListUsers(ctx, reqInvalidAdminToken)
	if err != nil {
		return result, err
	}
	if resInvalidAdminToken.StatusCode == 400 {
		result.WithInvalidAdminToken = true
	}

	reqValidWithoutQueryParams := sdk.NewListUsersReq(VALID_ADMIN_TOKEN, VALID_API_KEY)
	resValidWithoutQueryParams, err := api.ListUsers(ctx, reqValidWithoutQueryParams)
	if err != nil {
		return result, err
	}
	if resValidWithoutQueryParams.StatusCode == 200 {
		result.ValidOperationWithoutQueryParams = true
	}

	reqValidWithQueryParams := sdk.NewListUsersReq(VALID_ADMIN_TOKEN, VALID_API_KEY).WithPageSize(&PAGE_SIZE)
	resValidWithQueryParams, err := api.ListUsers(ctx, reqValidWithQueryParams)
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

func testGetUser(ctx context.Context, api *sdk.TestingAPI) (GetUserResult, error) {
	var result GetUserResult
	invalidReq := sdk.NewGetUserReq("", VALID_API_KEY, VALID_SESSION_TOKEN)
	_, err := api.GetUser(ctx, invalidReq)
	if err != nil {
		if err.Reason == sdk.ReasonEncoding {
			result.WithoutPathParam = true
		} else {
			return result, err
		}
	}

	reqValid := sdk.NewGetUserReq("1", VALID_API_KEY, VALID_SESSION_TOKEN)
	resValid, err := api.GetUser(ctx, reqValid)
	if err != nil {
		return result, err
	}
	if resValid.StatusCode == 200 && resValid.Response200.Body.UserId == "1" {
		result.ValidOperation = true
	}
	return result, nil
}

type CreateUserResult struct {
	ValidOperationWithoutOptionalField bool
	ValidOperationWithOptionalField    bool
	ValidOperationWithArbitraryData    bool
}

func testCreateUser(ctx context.Context, api *sdk.TestingAPI) (CreateUserResult, error) {
	var result CreateUserResult

	reqOptionalFieldMissing := sdk.NewCreateUserReq(
		VALID_ADMIN_TOKEN,
		VALID_API_KEY,
		sdk.NewCreateUserRequestBody(
			"test@example.com",
			sdk.UserStatusACTIVE,
			"Test User",
		),
	)
	resOptionalFieldMissing, err := api.CreateUser(ctx, reqOptionalFieldMissing)
	if err != nil {
		return result, err
	}
	if resOptionalFieldMissing.StatusCode == 201 && resOptionalFieldMissing.Response201.Body.User.Email == "test@example.com" {
		if resOptionalFieldMissing.Response201.Body.Status == sdk.UserStatusACTIVE {
			result.ValidOperationWithoutOptionalField = true
		}
	}

	optionalUserStatus := sdk.UserStatusINACTIVE_USER

	reqOptionalFieldPresent := sdk.NewCreateUserReq(
		VALID_ADMIN_TOKEN,
		VALID_API_KEY,
		sdk.NewCreateUserRequestBody(
			"test@example.com",
			optionalUserStatus,
			"Test User",
		).WithAge(AGE).WithOptionalStatus(optionalUserStatus),
	)
	resOptionalFieldPresent, err := api.CreateUser(ctx, reqOptionalFieldPresent)
	if err != nil {
		return result, err
	}
	if resOptionalFieldPresent.StatusCode == 201 && resOptionalFieldPresent.Response201.Body.User.Age != nil && *resOptionalFieldPresent.Response201.Body.User.Age == AGE {
		// Check if the optional status is set correctly
		if resOptionalFieldPresent.Response201.Body.OptionalStatus != nil && *resOptionalFieldPresent.Response201.Body.OptionalStatus == optionalUserStatus {
			result.ValidOperationWithOptionalField = true
		}
	}

	arbitraryDataSent := map[string]any{
		"key1": "value1",
		"key2": 123.0,
		"key3": true,
	}

	reqWithArbitraryData := sdk.NewCreateUserReq(
		VALID_ADMIN_TOKEN,
		VALID_API_KEY,
		sdk.NewCreateUserRequestBody(
			"test@example.com",
			sdk.UserStatusINACTIVE_USER,
			"Test User",
		).WithArbitraryData(arbitraryDataSent),
	)
	resWithArbitraryData, err := api.CreateUser(ctx, reqWithArbitraryData)
	if err != nil {
		return result, err
	}
	if resWithArbitraryData.StatusCode == 201 && resWithArbitraryData.Response201.Body.User.Email == "test@example.com" {
		if resWithArbitraryData.Response201.Body.ArbitraryData != nil && maps.Equal(arbitraryDataSent, *resWithArbitraryData.Response201.Body.ArbitraryData) {
			if resWithArbitraryData.Response201.Body.OptionalStatus == nil && resWithArbitraryData.Response201.Body.Status == sdk.UserStatusINACTIVE_USER {
				result.ValidOperationWithArbitraryData = true
			}
		}
	}

	return result, nil
}

type WhoAmIResult struct {
	ValidRawBody bool
}

func testWhoAmI(ctx context.Context, api *sdk.TestingAPI) (WhoAmIResult, error) {
	var result WhoAmIResult

	userId := "test@example.com"
	req := sdk.NewWhoAmIReq(VALID_API_KEY, VALID_SESSION_TOKEN)
	res, err := api.WhoAmI(ctx, req, strings.NewReader(userId))
	if err != nil {
		return result, err
	}

	if res.StatusCode == 200 {
		if bytes, err := io.ReadAll(res.Response200.RawBody.Body); err != nil {
			return result, err
		} else {
			if userId == string(bytes) {
				result.ValidRawBody = true
			}
		}
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
