package main

import (
	"net/http"

	_ "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// UserPostStruct holds all values of an incoming POST request
type UserPostStruct struct {
	User struct {
		Email string `json:"email"`
	} `json:"user"`
}

// PostDoc returns the description of this API endpoint
func (r *UserResource) PostDoc() string {
	return "create a new user invitation"
}

// PostParams returns the parameters supported by this API endpoint
func (r *UserResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *UserResource) Post(context smolder.APIContext, request *restful.Request, response *restful.Response, auth interface{}) {
	if auth == nil || auth.(DbUser).ID != 1 {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"UserResource POST"))
		return
	}

	ups := UserPostStruct{}
	err := request.ReadEntity(&ups)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse POST data",
			"UserResource POST"))
		return
	}

	user := DbUser{
		Username: ups.User.Email,
		Email:    ups.User.Email,
	}
	err = user.Save(context.(*PollyContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't create user",
			"UserResource POST"))
		return
	}

	sendInvitation(&user)

	resp := UserResponse{}
	resp.Init(context)
	resp.AddUser(&user)
	resp.Send(response)
}
