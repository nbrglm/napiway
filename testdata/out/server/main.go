package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nbrglm/napiway/testdata/out/server/api"
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

	mux.HandleFunc(api.ListUsersRequestRoutePath, handleListUsers)

	mux.HandleFunc(api.GetUserRequestRoutePath, handleGetUser)

	mux.HandleFunc(api.CreateUserRequestRoutePath, handleCreateUser)

	mux.HandleFunc(api.LogoutUserRequestRoutePath, handleLogoutUser)

	mux.HandleFunc(api.HealthCheckRequestRoutePath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != api.HealthCheckRequestHTTPMethod {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		api.WriteHealthCheck200Response(w, api.HealthCheck200Response{
			Body: api.HealthCheck200ResponseBody{
				Status: "ok",
			},
		})
	})

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		stdErr(true, "Failed to start server OR Exiting Server: %v\n", err)
	}
}

func handleListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.ListUsersRequestHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.NewListUsersRequest(w, r)
	if err != nil {
		debugMsg := err.Error()
		api.WriteListUsers400Response(w, api.ListUsers400Response{
			Body: api.ListUsers400ResponseBody{
				Error: api.ListUsers400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Invalid request",
				},
			},
		})
		return
	}

	if *req.Auth.APIKey != "valid" {
		debugMsg := "Invalid API key"
		api.WriteListUsers400Response(w, api.ListUsers400Response{
			Body: api.ListUsers400ResponseBody{
				Error: api.ListUsers400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Unauthorized",
				},
			},
		})
		return
	}

	if *req.Auth.AdminToken != "valid" {
		debugMsg := "Invalid Admin token"
		api.WriteListUsers400Response(w, api.ListUsers400Response{
			Body: api.ListUsers400ResponseBody{
				Error: api.ListUsers400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Unauthorized",
				},
			},
		})
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

	respUsers := make([]api.ListUsers200ResponseBodyUsersItem, 0, endIndex-startIndex)
	for _, user := range users[startIndex:endIndex] {
		respUsers = append(respUsers, api.ListUsers200ResponseBodyUsersItem{
			UserId:   user.ID,
			UserName: user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
			Age:      user.Age,
		})
	}

	api.WriteListUsers200Response(w, api.ListUsers200Response{
		Body: api.ListUsers200ResponseBody{
			Users:      respUsers,
			PageNumber: float64(pageNumber),
			PageSize:   float64(pageSize),
			TotalCount: float64(len(users)),
		},
	})
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.GetUserRequestHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.NewGetUserRequest(w, r)
	if err != nil {
		debugMsg := err.Error()
		api.WriteGetUser400Response(w, api.GetUser400Response{
			Body: api.GetUser400ResponseBody{
				Error: api.GetUser400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Invalid request",
				},
			},
		})
		return
	}

	if *req.Auth.APIKey != "valid" {
		debugMsg := "Invalid API key"
		api.WriteGetUser400Response(w, api.GetUser400Response{
			Body: api.GetUser400ResponseBody{
				Error: api.GetUser400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Unauthorized",
				},
			},
		})
		return
	}

	if *req.Auth.SessionToken != "valid" {
		debugMsg := "Invalid Session token"
		api.WriteGetUser400Response(w, api.GetUser400Response{
			Body: api.GetUser400ResponseBody{
				Error: api.GetUser400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Unauthorized",
				},
			},
		})
		return
	}

	user := new(api.GetUser200ResponseBodyUser)
	for _, u := range users {
		if u.ID == req.UserId {
			user = &api.GetUser200ResponseBodyUser{
				UserId:   u.ID,
				UserName: u.Name,
				Email:    u.Email,
				IsActive: u.IsActive,
				Age:      u.Age,
			}
			break
		}
	}

	if user == nil {
		debugMsg := "User not found"
		api.WriteGetUser404Response(w, api.GetUser404Response{
			Body: api.GetUser404ResponseBody{
				Error: api.GetUser404ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "User not found",
				},
			},
		})
		return
	}

	api.WriteGetUser200Response(w, api.GetUser200Response{
		Body: api.GetUser200ResponseBody{
			User: *user,
		},
	})
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.CreateUserRequestHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.NewCreateUserRequest(w, r)

	if err != nil {
		debugMsg := err.Error()
		api.WriteCreateUser400Response(w, api.CreateUser400Response{
			Body: api.CreateUser400ResponseBody{
				Error: api.CreateUser400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Invalid request",
				},
			},
		})
		return
	}

	if *req.Auth.APIKey != "valid" {
		debugMsg := "Invalid API key"
		api.WriteCreateUser400Response(w, api.CreateUser400Response{
			Body: api.CreateUser400ResponseBody{
				Error: api.CreateUser400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Unauthorized",
				},
			},
		})
		return
	}

	if *req.Auth.AdminToken != "valid" {
		debugMsg := "Invalid Admin token"
		api.WriteCreateUser400Response(w, api.CreateUser400Response{
			Body: api.CreateUser400ResponseBody{
				Error: api.CreateUser400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Unauthorized",
				},
			},
		})
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
	api.WriteCreateUser201Response(w, api.CreateUser201Response{
		Body: api.CreateUser201ResponseBody{
			User: api.CreateUser201ResponseBodyUser{
				UserId:   user.ID,
				UserName: user.Name,
				Email:    user.Email,
				IsActive: user.IsActive,
				Age:      user.Age,
			},
		},
	})
}

func handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != api.LogoutUserRequestHTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := api.NewLogoutUserRequest(w, r)
	if err != nil {
		debugMsg := err.Error()
		api.WriteLogoutUser400Response(w, api.LogoutUser400Response{
			Body: api.LogoutUser400ResponseBody{
				Error: api.LogoutUser400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Invalid request",
				},
			},
		})
		return
	}

	if *req.Auth.APIKey == "valid" {
		if req.Auth.RefreshToken != nil && *req.Auth.RefreshToken == "valid" {
			api.WriteLogoutUser200Response(w, api.LogoutUser200Response{
				Body: api.LogoutUser200ResponseBody{
					Message: "RefreshToken Logout successful",
				},
			})
		} else if req.Auth.SessionToken != nil && *req.Auth.SessionToken == "valid" {
			api.WriteLogoutUser200Response(w, api.LogoutUser200Response{
				Body: api.LogoutUser200ResponseBody{
					Message: "SessionToken Logout successful",
				},
			})
		} else {
			debugMsg := "Invalid or missing token"
			api.WriteLogoutUser400Response(w, api.LogoutUser400Response{
				Body: api.LogoutUser400ResponseBody{
					Error: api.LogoutUser400ResponseBodyError{
						DebugMessage: &debugMsg,
						ErrorMessage: "Unauthorized",
					},
				},
			})
		}
	} else {
		debugMsg := "Invalid API key"
		api.WriteLogoutUser400Response(w, api.LogoutUser400Response{
			Body: api.LogoutUser400ResponseBody{
				Error: api.LogoutUser400ResponseBodyError{
					DebugMessage: &debugMsg,
					ErrorMessage: "Unauthorized",
				},
			},
		})
	}
}

func stdErr(exit bool, format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	if exit {
		os.Exit(1)
	}
}
