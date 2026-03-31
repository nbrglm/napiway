package go_sdk

import (
	"net/http"
)

const (
	WhoAmIReqHTTPMethod = "POST"
	WhoAmIReqRoutePath  = "/users/whoami"
)

// Get information about the currently authenticated user. Provide the request body as just a string which is set to the user id.
type WhoAmIReq struct {

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

// Successful response containing information about the currently authenticated user. Body is just a string with the user id provided in request body.
type WhoAmI200 struct {

	// Source: header parameter "X-RateLimit-Remaining"
	//

	// The number of remaining requests allowed in the current rate limit window.
	//
	// Required
	RateLimitRemaining int64

	// Raw response body. The HTTP response will be returned directly for this response, and it will be the responsibility of the caller to read/close the response body.
	RawBody *http.Response
}

// WhoAmI400 represents a response with no headers and no response body
// Invalid Request

// NewWhoAmIReq creates a new instance of WhoAmIReq with required fields as parameters
func NewWhoAmIReq(

	APIKeyAuth string,

	SessionTokenAuth string,

) *WhoAmIReq {
	return &WhoAmIReq{

		APIKeyAuth: APIKeyAuth,

		SessionTokenAuth: SessionTokenAuth,
	}
}

// NOTE: RawBody is true, so request body will not be handled.

// ParseWhoAmI200 creates a new instance of WhoAmI200 by parsing a map[string]any
func ParseWhoAmI200(resp *http.Response) (*WhoAmI200, error) {
	result := new(WhoAmI200)

	headerRateLimitRemaining, err := parseint64Param(resp.Header.Get("X-RateLimit-Remaining"), "header: X-RateLimit-Remaining", true)
	if err != nil {
		return nil, err
	}

	result.RateLimitRemaining = *headerRateLimitRemaining

	result.RawBody = resp

	return result, nil
}
