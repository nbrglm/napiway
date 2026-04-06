package api

import (
	"fmt"
	"net/http"
	"strings"
)

// GetAdminToken extracts the AdminToken Authentication (header: "X-App-Admin-Token") from the request and returns it as a string.
func GetAdminToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("X-App-Admin-Token")
	if authHeader == "" {
		return "", fmt.Errorf("missing auth header: %s", "AdminToken")
	}
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", fmt.Errorf("empty auth header: %s", "AdminToken")
	}
	return authHeader, nil
}

// GetAPIKey extracts the APIKey Authentication (header: "X-App-API-Key") from the request and returns it as a string.
func GetAPIKey(r *http.Request) (string, error) {
	authHeader := r.Header.Get("X-App-API-Key")
	if authHeader == "" {
		return "", fmt.Errorf("missing auth header: %s", "APIKey")
	}
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", fmt.Errorf("empty auth header: %s", "APIKey")
	}
	return authHeader, nil
}

// GetRefreshToken extracts the RefreshToken Authentication (header: "X-App-Refresh-Token") from the request and returns it as a string.
func GetRefreshToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("X-App-Refresh-Token")
	if authHeader == "" {
		return "", fmt.Errorf("missing auth header: %s", "RefreshToken")
	}
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", fmt.Errorf("empty auth header: %s", "RefreshToken")
	}
	return authHeader, nil
}

// GetSessionToken extracts the SessionToken Authentication (header: "X-App-Session-Token") from the request and returns it as a string.
func GetSessionToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("X-App-Session-Token")
	if authHeader == "" {
		return "", fmt.Errorf("missing auth header: %s", "SessionToken")
	}
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", fmt.Errorf("empty auth header: %s", "SessionToken")
	}
	return authHeader, nil
}

type CreateUserRequestBody struct {

	// The age of the user to be created.
	//
	// Optional
	//
	Age *int64 `json:"Age,omitempty"`

	// An object that can contain any arbitrary data related to the user list response. Just for testing freeForm object support in the generator. This is returned by the server in the response body.
	//
	// Optional
	//
	ArbitraryData *map[string]any `json:"ArbitraryData,omitempty"`

	// The email address of the user to be created.
	//
	// Required
	//
	// Must be non-empty
	Email string `json:"Email"`

	// An optional status of the user to be created. Just for testing optional enum support in the generator.
	//
	// Optional
	//
	OptionalStatus *UserStatus `json:"OptionalStatus,omitempty"`

	// The status of the user to be created.
	//
	// Required
	//
	Status UserStatus `json:"Status"`

	// The name of the user to be created.
	//
	// Required
	//
	// Must be non-empty
	UserName string `json:"UserName"`
}

// NewCreateUserRequestBody creates a new instance of CreateUserRequestBody with required fields as parameters
func NewCreateUserRequestBody(

	Email string,

	Status UserStatus,

	UserName string,

) *CreateUserRequestBody {
	return &CreateUserRequestBody{

		Email: Email,

		Status: Status,

		UserName: UserName,
	}
}

// WithAge sets the optional field Age and returns the modified CreateUserRequestBody instance
func (o *CreateUserRequestBody) WithAge(value *int64) *CreateUserRequestBody {
	o.Age = value
	return o
}

// WithArbitraryData sets the optional field ArbitraryData and returns the modified CreateUserRequestBody instance
func (o *CreateUserRequestBody) WithArbitraryData(value *map[string]any) *CreateUserRequestBody {
	o.ArbitraryData = value
	return o
}

// WithOptionalStatus sets the optional field OptionalStatus and returns the modified CreateUserRequestBody instance
func (o *CreateUserRequestBody) WithOptionalStatus(value *UserStatus) *CreateUserRequestBody {
	o.OptionalStatus = value
	return o
}

func ParseCreateUserRequestBody(data map[string]any) (*CreateUserRequestBody, error) {
	body := new(CreateUserRequestBody)

	valAge, ok := data["Age"]
	if !ok {

		// skip, leave as zero value

	} else {

		var valAgeTyped int64
		// JSON numbers are float64 by default, so we need to handle that case
		switch v := valAge.(type) {
		case float64:
			valAgeTyped = int64(v)
		case int64:
			valAgeTyped = v
		default:
			return body, fmt.Errorf("field 'Age' has incorrect type")
		}

		body.Age = &valAgeTyped

	}

	valArbitraryData, ok := data["ArbitraryData"]
	if !ok {

		// skip, leave as zero value

	} else {
		valArbitraryDataTyped, ok := valArbitraryData.(map[string]any)
		if !ok {
			return body, fmt.Errorf("field 'ArbitraryData' has incorrect type")
		}

		body.ArbitraryData = &valArbitraryDataTyped
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

	valOptionalStatus, ok := data["OptionalStatus"]
	if !ok {

		// skip, leave as zero value

	} else {

		valOptionalStatusStr, ok := valOptionalStatus.(string)
		if !ok {
			return body, fmt.Errorf("field 'OptionalStatus' has incorrect type")
		}
		valOptionalStatusTyped, err := ParseUserStatus(valOptionalStatusStr)
		if err != nil {
			return body, fmt.Errorf("field 'OptionalStatus' is invalid: %w", err)
		}

		body.OptionalStatus = valOptionalStatusTyped

	}

	valStatus, ok := data["Status"]
	if !ok {

		return body, fmt.Errorf("missing required field 'Status'")

	} else {

		valStatusStr, ok := valStatus.(string)
		if !ok {
			return body, fmt.Errorf("field 'Status' has incorrect type")
		}
		valStatusTyped, err := ParseUserStatus(valStatusStr)
		if err != nil {
			return body, fmt.Errorf("field 'Status' is invalid: %w", err)
		}

		body.Status = *valStatusTyped

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

type CreateUserResponseBody struct {

	// An object that can contain any arbitrary data related to the user list response. Just for testing freeForm object support in the generator.
	//
	// Optional
	//
	ArbitraryData *map[string]any `json:"ArbitraryData,omitempty"`

	// An optional status of the user. Just for testing optional enum support in the generator.
	//
	// Optional
	//
	OptionalStatus *UserStatus `json:"OptionalStatus,omitempty"`

	// The status of the user.
	//
	// Required
	//
	Status UserStatus `json:"Status"`

	// The created user information.
	//
	// Required
	//
	User *User `json:"User"`
}

// NewCreateUserResponseBody creates a new instance of CreateUserResponseBody with required fields as parameters
func NewCreateUserResponseBody(

	Status UserStatus,

	User *User,

) *CreateUserResponseBody {
	return &CreateUserResponseBody{

		Status: Status,

		User: User,
	}
}

// WithArbitraryData sets the optional field ArbitraryData and returns the modified CreateUserResponseBody instance
func (o *CreateUserResponseBody) WithArbitraryData(value *map[string]any) *CreateUserResponseBody {
	o.ArbitraryData = value
	return o
}

// WithOptionalStatus sets the optional field OptionalStatus and returns the modified CreateUserResponseBody instance
func (o *CreateUserResponseBody) WithOptionalStatus(value *UserStatus) *CreateUserResponseBody {
	o.OptionalStatus = value
	return o
}

func ParseCreateUserResponseBody(data map[string]any) (*CreateUserResponseBody, error) {
	body := new(CreateUserResponseBody)

	valArbitraryData, ok := data["ArbitraryData"]
	if !ok {

		// skip, leave as zero value

	} else {
		valArbitraryDataTyped, ok := valArbitraryData.(map[string]any)
		if !ok {
			return body, fmt.Errorf("field 'ArbitraryData' has incorrect type")
		}

		body.ArbitraryData = &valArbitraryDataTyped
	}

	valOptionalStatus, ok := data["OptionalStatus"]
	if !ok {

		// skip, leave as zero value

	} else {

		valOptionalStatusStr, ok := valOptionalStatus.(string)
		if !ok {
			return body, fmt.Errorf("field 'OptionalStatus' has incorrect type")
		}
		valOptionalStatusTyped, err := ParseUserStatus(valOptionalStatusStr)
		if err != nil {
			return body, fmt.Errorf("field 'OptionalStatus' is invalid: %w", err)
		}

		body.OptionalStatus = valOptionalStatusTyped

	}

	valStatus, ok := data["Status"]
	if !ok {

		return body, fmt.Errorf("missing required field 'Status'")

	} else {

		valStatusStr, ok := valStatus.(string)
		if !ok {
			return body, fmt.Errorf("field 'Status' has incorrect type")
		}
		valStatusTyped, err := ParseUserStatus(valStatusStr)
		if err != nil {
			return body, fmt.Errorf("field 'Status' is invalid: %w", err)
		}

		body.Status = *valStatusTyped

	}

	valUser, ok := data["User"]
	if !ok {

		return body, fmt.Errorf("missing required field 'User'")

	} else {

		valUserMap, ok := valUser.(map[string]any)
		if !ok {
			return body, fmt.Errorf("field 'User' has incorrect type")
		}
		valUserTyped, err := ParseUser(valUserMap)
		if err != nil {
			return body, fmt.Errorf("field 'User' is invalid: %w", err)
		}

		body.User = valUserTyped

	}

	return body, nil
}

type ErrorResponse struct {

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

// NewErrorResponse creates a new instance of ErrorResponse with required fields as parameters
func NewErrorResponse(

	ErrorMessage string,

) *ErrorResponse {
	return &ErrorResponse{

		ErrorMessage: ErrorMessage,
	}
}

// WithDebugMessage sets the optional field DebugMessage and returns the modified ErrorResponse instance
func (o *ErrorResponse) WithDebugMessage(value *string) *ErrorResponse {
	o.DebugMessage = value
	return o
}

func ParseErrorResponse(data map[string]any) (*ErrorResponse, error) {
	body := new(ErrorResponse)

	valDebugMessage, ok := data["DebugMessage"]
	if !ok {

		// skip, leave as zero value

	} else {

		valDebugMessageTyped, ok := valDebugMessage.(string)
		if !ok {
			return body, fmt.Errorf("field 'DebugMessage' has incorrect type")
		}

		body.DebugMessage = &valDebugMessageTyped

	}

	valErrorMessage, ok := data["ErrorMessage"]
	if !ok {

		return body, fmt.Errorf("missing required field 'ErrorMessage'")

	} else {

		valErrorMessageTyped, ok := valErrorMessage.(string)
		if !ok {
			return body, fmt.Errorf("field 'ErrorMessage' has incorrect type")
		}

		valErrorMessageTyped = strings.TrimSpace(valErrorMessageTyped)
		if len(valErrorMessageTyped) == 0 {
			return body, fmt.Errorf("field 'ErrorMessage' must be non-empty")
		}

		body.ErrorMessage = valErrorMessageTyped

	}

	return body, nil
}

type HealthCheckResponseBody struct {

	// The health status of the API, typically "OK".
	//
	// Required
	//
	// Must be non-empty
	Status string `json:"Status"`
}

// NewHealthCheckResponseBody creates a new instance of HealthCheckResponseBody with required fields as parameters
func NewHealthCheckResponseBody(

	Status string,

) *HealthCheckResponseBody {
	return &HealthCheckResponseBody{

		Status: Status,
	}
}

func ParseHealthCheckResponseBody(data map[string]any) (*HealthCheckResponseBody, error) {
	body := new(HealthCheckResponseBody)

	valStatus, ok := data["Status"]
	if !ok {

		return body, fmt.Errorf("missing required field 'Status'")

	} else {

		valStatusTyped, ok := valStatus.(string)
		if !ok {
			return body, fmt.Errorf("field 'Status' has incorrect type")
		}

		valStatusTyped = strings.TrimSpace(valStatusTyped)
		if len(valStatusTyped) == 0 {
			return body, fmt.Errorf("field 'Status' must be non-empty")
		}

		body.Status = valStatusTyped

	}

	return body, nil
}

type ListUsersResponseBody struct {

	// The current page number.
	//
	// Required
	//
	PageNumber int64 `json:"PageNumber"`

	// The number of items per page.
	//
	// Required
	//
	PageSize int64 `json:"PageSize"`

	// The total number of users available.
	//
	// Required
	//
	TotalCount int64 `json:"TotalCount"`

	// Required
	//
	Users []User `json:"Users"`
}

// NewListUsersResponseBody creates a new instance of ListUsersResponseBody with required fields as parameters
func NewListUsersResponseBody(

	PageNumber int64,

	PageSize int64,

	TotalCount int64,

	Users []User,

) *ListUsersResponseBody {
	return &ListUsersResponseBody{

		PageNumber: PageNumber,

		PageSize: PageSize,

		TotalCount: TotalCount,

		Users: Users,
	}
}

func ParseListUsersResponseBody(data map[string]any) (*ListUsersResponseBody, error) {
	body := new(ListUsersResponseBody)

	valPageNumber, ok := data["PageNumber"]
	if !ok {

		return body, fmt.Errorf("missing required field 'PageNumber'")

	} else {

		var valPageNumberTyped int64
		// JSON numbers are float64 by default, so we need to handle that case
		switch v := valPageNumber.(type) {
		case float64:
			valPageNumberTyped = int64(v)
		case int64:
			valPageNumberTyped = v
		default:
			return body, fmt.Errorf("field 'PageNumber' has incorrect type")
		}

		body.PageNumber = valPageNumberTyped

	}

	valPageSize, ok := data["PageSize"]
	if !ok {

		return body, fmt.Errorf("missing required field 'PageSize'")

	} else {

		var valPageSizeTyped int64
		// JSON numbers are float64 by default, so we need to handle that case
		switch v := valPageSize.(type) {
		case float64:
			valPageSizeTyped = int64(v)
		case int64:
			valPageSizeTyped = v
		default:
			return body, fmt.Errorf("field 'PageSize' has incorrect type")
		}

		body.PageSize = valPageSizeTyped

	}

	valTotalCount, ok := data["TotalCount"]
	if !ok {

		return body, fmt.Errorf("missing required field 'TotalCount'")

	} else {

		var valTotalCountTyped int64
		// JSON numbers are float64 by default, so we need to handle that case
		switch v := valTotalCount.(type) {
		case float64:
			valTotalCountTyped = int64(v)
		case int64:
			valTotalCountTyped = v
		default:
			return body, fmt.Errorf("field 'TotalCount' has incorrect type")
		}

		body.TotalCount = valTotalCountTyped

	}

	valUsers, ok := data["Users"]
	if !ok {

		return body, fmt.Errorf("missing required field 'Users'")

	} else {

		valUsersSlice, ok := valUsers.([]any)
		if !ok {
			return body, fmt.Errorf("field 'Users' has incorrect type")
		}

		if len(valUsersSlice) == 0 {
			return body, fmt.Errorf("field 'Users' must be non-empty")
		}

		valUsersTyped := make([]User, 0, len(valUsersSlice))

		for idx, item := range valUsersSlice {

			itemMap, ok := item.(map[string]any)
			if !ok {
				return body, fmt.Errorf("element %d of field 'Users' has incorrect type", idx)
			}
			validatedItem, err := ParseUser(itemMap)
			if err != nil {
				return body, fmt.Errorf("element %d of field 'Users' is invalid: %w", idx, err)
			}

			valUsersTyped = append(valUsersTyped, *validatedItem)
		}

		body.Users = valUsersTyped

	}

	return body, nil
}

type LogoutUserResponseBody struct {

	// A message confirming successful logout.
	//
	// Required
	//
	// Must be non-empty
	Message string `json:"Message"`
}

// NewLogoutUserResponseBody creates a new instance of LogoutUserResponseBody with required fields as parameters
func NewLogoutUserResponseBody(

	Message string,

) *LogoutUserResponseBody {
	return &LogoutUserResponseBody{

		Message: Message,
	}
}

func ParseLogoutUserResponseBody(data map[string]any) (*LogoutUserResponseBody, error) {
	body := new(LogoutUserResponseBody)

	valMessage, ok := data["Message"]
	if !ok {

		return body, fmt.Errorf("missing required field 'Message'")

	} else {

		valMessageTyped, ok := valMessage.(string)
		if !ok {
			return body, fmt.Errorf("field 'Message' has incorrect type")
		}

		valMessageTyped = strings.TrimSpace(valMessageTyped)
		if len(valMessageTyped) == 0 {
			return body, fmt.Errorf("field 'Message' must be non-empty")
		}

		body.Message = valMessageTyped

	}

	return body, nil
}

type User struct {

	// The age of the user.
	//
	// Optional
	//
	Age *int64 `json:"Age,omitempty"`

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

// NewUser creates a new instance of User with required fields as parameters
func NewUser(

	Email string,

	IsActive bool,

	UserId string,

	UserName string,

) *User {
	return &User{

		Email: Email,

		IsActive: IsActive,

		UserId: UserId,

		UserName: UserName,
	}
}

// WithAge sets the optional field Age and returns the modified User instance
func (o *User) WithAge(value *int64) *User {
	o.Age = value
	return o
}

func ParseUser(data map[string]any) (*User, error) {
	body := new(User)

	valAge, ok := data["Age"]
	if !ok {

		// skip, leave as zero value

	} else {

		var valAgeTyped int64
		// JSON numbers are float64 by default, so we need to handle that case
		switch v := valAge.(type) {
		case float64:
			valAgeTyped = int64(v)
		case int64:
			valAgeTyped = v
		default:
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

	valIsActive, ok := data["IsActive"]
	if !ok {

		return body, fmt.Errorf("missing required field 'IsActive'")

	} else {

		valIsActiveTyped, ok := valIsActive.(bool)
		if !ok {
			return body, fmt.Errorf("field 'IsActive' has incorrect type")
		}

		body.IsActive = valIsActiveTyped

	}

	valUserId, ok := data["UserId"]
	if !ok {

		return body, fmt.Errorf("missing required field 'UserId'")

	} else {

		valUserIdTyped, ok := valUserId.(string)
		if !ok {
			return body, fmt.Errorf("field 'UserId' has incorrect type")
		}

		valUserIdTyped = strings.TrimSpace(valUserIdTyped)
		if len(valUserIdTyped) == 0 {
			return body, fmt.Errorf("field 'UserId' must be non-empty")
		}

		body.UserId = valUserIdTyped

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

type UserStatus string

const (
	UserStatusACTIVE UserStatus = "ACTIVE"

	UserStatusINACTIVE_USER UserStatus = "INACTIVE_USER"
)

// NewUserStatus isn't required since enums are just strings.

// ParseUserStatus parses a string into a UserStatus value, returning an error if the input is not a valid enum value.
func ParseUserStatus(data string) (*UserStatus, error) {
	switch data {

	case "ACTIVE":
		var enumValue UserStatus = UserStatusACTIVE
		return &enumValue, nil

	case "INACTIVE_USER":
		var enumValue UserStatus = UserStatusINACTIVE_USER
		return &enumValue, nil

	default:
		return nil, fmt.Errorf("invalid value for UserStatus: %s", data)
	}
}
