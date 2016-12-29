package main

import (
	_ "github.com/Sirupsen/logrus"
	"github.com/muesli/smolder"
)

// UserResponse is the common response to 'user' requests
type UserResponse struct {
	smolder.Response

	Users []userInfoResponse `json:"users,omitempty"`
	users []DbUser
}

type userInfoResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	About     string `json:"about"`
	Activated bool   `json:"activated"`
	URL       string `json:"url"`
}

// Init a new response
func (r *UserResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Users = []userInfoResponse{}
}

// AddUser adds a user to the response
func (r *UserResponse) AddUser(user *DbUser) {
	r.users = append(r.users, *user)
	r.Users = append(r.Users, prepareUserResponse(r.Context, user))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *UserResponse) EmptyResponse() interface{} {
	if len(r.users) == 0 {
		var out struct {
			Users interface{} `json:"users"`
		}
		out.Users = []userInfoResponse{}
		return out
	}
	return nil
}

func prepareUserResponse(context smolder.APIContext, user *DbUser) userInfoResponse {
	resp := userInfoResponse{
		ID:        user.ID,
		Email:     user.Email,
		About:     user.About,
		Activated: user.Activated,
		URL:       buildURL(*user),
	}

	return resp
}
