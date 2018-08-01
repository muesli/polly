package users

import (
	"net/http"

	"github.com/muesli/polly/api/db"
	"github.com/muesli/polly/api/utils"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// UserPostStruct holds all values of an incoming POST request
type UserPostStruct struct {
	User struct {
		Email string `json:"email"`
		About string `json:"about"`
	} `json:"user"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *UserResource) PostAuthRequired() bool {
	return true
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
func (r *UserResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	auth, err := context.Authentication(request)
	if err != nil || auth.(db.User).ID != 1 {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"UserResource POST"))
		return
	}

	ups := UserPostStruct{}
	err = request.ReadEntity(&ups)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse POST data",
			"UserResource POST"))
		return
	}

	if ups.User.About == "" {
		ups.User.About = ups.User.Email
	}

	user := db.User{
		Username: ups.User.Email,
		Email:    ups.User.Email,
		About:    ups.User.About,
	}
	err = user.Save(context.(*db.PollyContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't create user",
			"UserResource POST"))
		return
	}

	utils.SendInvitation(&user)

	resp := UserResponse{}
	resp.Init(context)
	resp.AddUser(&user)
	resp.Send(response)
}
