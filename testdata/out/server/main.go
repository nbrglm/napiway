package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/nbrglm/napiway/testdata/out/server/api"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
	Age      *int64 `json:"age"`
}

var age1 = int64(28)

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

func main() {
	exitSignal := make(chan os.Signal, 1)
	go func() {
		<-exitSignal
		stdErr(false, "Exiting server...\n")
		os.Exit(0)
	}()
	args := os.Args
	fmt.Printf("Starting server with args: %v\n", args)
	if len(args) < 2 || args[1] == "" {
		stdErr(true, "Please provide the server address as the first argument")
		return
	}
	serverAddr := args[1]

	mux := http.NewServeMux()

	mux.HandleFunc(api.ListUsersReqRoutePath, handleListUsers)

	mux.HandleFunc(api.GetUserReqRoutePath, handleGetUser)

	mux.HandleFunc(api.CreateUserReqRoutePath, handleCreateUser)

	mux.HandleFunc(api.LogoutUserReqRoutePath, handleLogoutUser)

	mux.HandleFunc(api.WhoAmIReqRoutePath, handleWhoAmI)

	mux.HandleFunc(api.HealthCheckReqRoutePath, func(w http.ResponseWriter, r *http.Request) {
		req, _ := api.ParseHealthCheckReq(w, r)
		if r.Method != api.HealthCheckReqHTTPMethod {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		req.Write200(
			w,
			api.NewHealthCheck200(
				api.NewHealthCheckResponseBody("ok"),
			),
		)
	})

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		stdErr(true, "Failed to start server OR Exiting Server: %v\n", err)
	}
}

func handleListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.ListUsersReqHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.ParseListUsersReq(w, r)
	if err != nil {
		debugMsg := err.Error()
		req.Write400(w, api.NewListUsers400(
			api.NewErrorResponse(
				debugMsg,
			),
		))
		return
	}

	if req.APIKeyAuth != "valid" {
		debugMsg := "Invalid API key"
		req.Write400(w, api.NewListUsers400(
			api.NewErrorResponse(
				debugMsg,
			),
		))
		return
	}

	if req.AdminTokenAuth != "valid" {
		debugMsg := "Invalid Admin token"
		req.Write400(w, api.NewListUsers400(
			api.NewErrorResponse(
				debugMsg,
			),
		))
		return
	}

	pageNumber := 0
	if req.PageNumber != nil && *req.PageNumber >= 0 {
		pageNumber = int(*req.PageNumber)
	}

	pageSize := 10
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = int(*req.PageSize)
	}

	startIndex := pageNumber * pageSize
	endIndex := startIndex + pageSize
	if startIndex > len(users) {
		startIndex = len(users)
	}
	if endIndex > len(users) {
		endIndex = len(users)
	}

	respUsers := make([]api.User, 0, endIndex-startIndex)
	for _, user := range users[startIndex:endIndex] {
		respUsers = append(respUsers, *mapToApiUser(user))
	}

	req.Write200(w, api.NewListUsers200(
		rand.Int64N(100),
		api.NewListUsersResponseBody(
			int64(pageNumber),
			int64(pageSize),
			int64(len(users)),
			respUsers,
		),
	))
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.GetUserReqHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.ParseGetUserReq(w, r)
	if err != nil {
		debugMsg := err.Error()
		fmt.Printf("Error parsing GetUser request: %s\n", debugMsg)
		req.Write400(
			w,
			api.NewGetUser400(
				api.NewErrorResponse("Invalid request"),
			),
		)
		return
	}

	if req.APIKeyAuth != "valid" {
		debugMsg := "Invalid API key"
		fmt.Printf("Error handling GetUser request: %s\n", debugMsg)
		req.Write400(
			w,
			api.NewGetUser400(
				api.NewErrorResponse("Unauthorized"),
			),
		)
		return
	}

	if req.SessionTokenAuth != "valid" {
		debugMsg := "Invalid Session token"
		fmt.Printf("Error handling GetUser request: %s\n", debugMsg)
		req.Write400(
			w,
			api.NewGetUser400(
				api.NewErrorResponse("Unauthorized"),
			),
		)
		return
	}

	user := new(api.User)
	for _, u := range users {
		if u.ID == req.UserId {
			user = mapToApiUser(u)
			break
		}
	}

	if user == nil {
		debugMsg := "User not found"
		fmt.Printf("Error handling GetUser request: %s\n", debugMsg)
		req.Write404(
			w,
			api.NewGetUser404(
				api.NewErrorResponse("User not found"),
			),
		)
		return
	}

	req.Write200(
		w,
		api.NewGetUser200(
			user,
		),
	)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.CreateUserReqHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.ParseCreateUserReq(w, r)

	if err != nil {
		debugMsg := err.Error()
		fmt.Printf("Error parsing CreateUser request: %s\n", debugMsg)
		req.Write400(w, api.NewCreateUser400(
			api.NewErrorResponse("Invalid request"),
		))
		return
	}

	if req.APIKeyAuth != "valid" {
		debugMsg := "Invalid API key"
		fmt.Printf("Error handling CreateUser request: %s\n", debugMsg)
		req.Write400(w, api.NewCreateUser400(
			api.NewErrorResponse(debugMsg),
		))
		return
	}

	if req.AdminTokenAuth != "valid" {
		debugMsg := "Invalid Admin token"
		fmt.Printf("Error handling CreateUser request: %s\n", debugMsg)
		req.Write400(w, api.NewCreateUser400(
			api.NewErrorResponse(debugMsg),
		))
		return
	}

	if req.Body.Age != nil && *req.Body.Age < 0 {
		debugMsg := "Age cannot be negative"
		fmt.Printf("Error handling CreateUser request: %s\n", debugMsg)
		req.Write400(w, api.NewCreateUser400(
			api.NewErrorResponse(debugMsg),
		))
		return
	}

	user := User{
		ID:       fmt.Sprintf("%d", len(users)+1),
		Name:     req.Body.UserName,
		Email:    req.Body.Email,
		IsActive: true,
		Age:      req.Body.Age,
	}

	users = append(users, user)
	respBody := api.NewCreateUserResponseBody(
		req.Body.Status,
		mapToApiUser(user),
	)
	if req.Body.ArbitraryData != nil {
		respBody = respBody.WithArbitraryData(*req.Body.ArbitraryData)
	}
	if req.Body.OptionalStatus != nil {
		respBody = respBody.WithOptionalStatus(*req.Body.OptionalStatus)
	}

	resp := api.NewCreateUser201(respBody)

	req.Write201(w, resp)
}

func handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.LogoutUserReqHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.ParseLogoutUserReq(w, r)
	if err != nil {
		debugMsg := err.Error()
		req.Write400(
			w,
			api.NewLogoutUser400(
				api.NewErrorResponse(debugMsg),
			),
		)
		return
	}

	if req.APIKeyAuth == "valid" {
		if req.RefreshTokenAuth != nil && *req.RefreshTokenAuth == "valid" {
			req.Write200(
				w,
				api.NewLogoutUser200(
					api.NewLogoutUserResponseBody(
						"RefreshToken Logout successful",
					),
				),
			)
		} else if req.SessionTokenAuth != nil && *req.SessionTokenAuth == "valid" {
			req.Write200(
				w,
				api.NewLogoutUser200(
					api.NewLogoutUserResponseBody(
						"SessionToken Logout successful",
					),
				),
			)
		} else {
			debugMsg := "Invalid or missing token"
			req.Write400(
				w,
				api.NewLogoutUser400(
					api.NewErrorResponse(debugMsg),
				),
			)
		}
	} else {
		debugMsg := "Invalid API key"
		req.Write400(
			w,
			api.NewLogoutUser400(
				api.NewErrorResponse(debugMsg),
			),
		)
	}
}

func handleWhoAmI(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.WhoAmIReqHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.ParseWhoAmIReq(w, r)
	if err != nil {
		debugMsg := err.Error()
		fmt.Printf("Error parsing WhoAmI request: %s\n", debugMsg)
		req.Write400(w)
		return
	}

	if req.APIKeyAuth != "valid" {
		debugMsg := "Invalid API key"
		fmt.Printf("Error handling WhoAmI request: %s\n", debugMsg)
		req.Write400(w)
		return
	}

	if req.SessionTokenAuth != "valid" {
		debugMsg := "Invalid Session token"
		fmt.Printf("Error handling WhoAmI request: %s\n", debugMsg)
		req.Write400(w)
		return
	}

	userIdBytes, err := io.ReadAll(r.Body)
	if err != nil {
		debugMsg := err.Error()
		fmt.Printf("Error reading WhoAmI request body: %s\n", debugMsg)
		req.Write400(w)
		return
	}

	req.Write200(w, api.NewWhoAmI200(10)) // To write the status code
	w.Write(userIdBytes)                  // Write the raw body
}

func stdErr(exit bool, format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	if exit {
		os.Exit(1)
	}
}

func mapToApiUser(user User) *api.User {
	u := api.NewUser(user.Email, user.IsActive, user.ID, user.Name)
	if user.Age != nil {
		u.WithAge(*user.Age)
	}
	return u
}
