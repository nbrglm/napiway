package api

import (
	"fmt"
	"net/http"
	"strings"
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

// ParseWhoAmIReq creates a new instance of WhoAmIReq by parsing the http.Request
func ParseWhoAmIReq(w http.ResponseWriter, r *http.Request) (*WhoAmIReq, error) {
	req := WhoAmIReq{}

	// Parse path parameters, if any

	// Parse query parameters, if any

	// Parse header parameters, if any

	// Required auth, if any

	valAPIKey := r.Header.Get("X-App-API-Key")
	valAPIKey = strings.TrimSpace(valAPIKey)
	if valAPIKey == "" {
		return nil, fmt.Errorf("missing required authentication: header X-App-API-Key")
	} else {
		req.APIKeyAuth = valAPIKey
	}

	valSessionToken := r.Header.Get("X-App-Session-Token")
	valSessionToken = strings.TrimSpace(valSessionToken)
	if valSessionToken == "" {
		return nil, fmt.Errorf("missing required authentication: header X-App-Session-Token")
	} else {
		req.SessionTokenAuth = valSessionToken
	}

	// Atleast one auth, if any

	// NOTE: RawBody is true, so request body will not be handled.

	return &req, nil
}

func NewWhoAmI200(

	RateLimitRemaining int64,

) *WhoAmI200 {
	return &WhoAmI200{

		RateLimitRemaining: RateLimitRemaining,
	}
}

// Write200 writes the WhoAmI200 response to the http.ResponseWriter
//
// RawBody is true, hence this function will only set the headers and write the status code, rest is to be done by the caller.
//
// Note: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
func (r *WhoAmIReq) Write200(w http.ResponseWriter, resp *WhoAmI200) error {
	// Set headers, if any

	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%v", resp.RateLimitRemaining))

	// Set status code and write the header as there are no body to write
	w.WriteHeader(200)
	return nil

}

// Write400 writes the WhoAmI400 response to the http.ResponseWriter
//
// # Invalid Request
//
// NOTE: THIS FUNCTION WILL CALL w.WriteHeader(), so ensure that no other calls to w.WriteHeader() are made before calling this function.
//
// Since there are no headers or body to write, this function will only set the status code in the response header.
func (r *WhoAmIReq) Write400(w http.ResponseWriter) error {
	// Set status code and write the header as there are no headers or body to write
	w.WriteHeader(400)
	return nil
}
