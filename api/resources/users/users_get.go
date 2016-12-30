package users

import (
	"net/http"
	"strconv"

	"github.com/muesli/polly/api/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *UserResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *UserResource) GetByIDsAuthRequired() bool {
	return true
}

// GetDoc returns the description of this API endpoint
func (r *UserResource) GetDoc() string {
	return "retrieve users"
}

// GetParams returns the parameters supported by this API endpoint
func (r *UserResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("token", "token of a user").DataType("string"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *UserResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := UserResponse{}
	resp.Init(context)

	for _, id := range ids {
		iid, err := strconv.Atoi(id)
		if err != nil {
			r.NotFound(request, response)
			return
		}
		user, err := context.(*db.PollyContext).GetUserByID(int64(iid))
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddUser(&user)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *UserResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := UserResponse{}
	resp.Init(context)

	token := params["token"]
	if len(token) > 0 {
		auth, err := context.(*db.PollyContext).GetUserByAccessToken(token[0])
		if auth == nil || err != nil {
			r.NotFound(request, response)
			return
		}
		user := auth.(db.User)

		resp.AddUser(&user)
	} else {
		auth, err := context.Authentication(request)
		if err != nil || auth == nil || auth.(db.User).ID != 1 {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				false,
				"Admin permission required for this operation",
				"UserResource GET"))
			return
		}

		users, err := context.(*db.PollyContext).LoadAllUsers()
		if err != nil {
			r.NotFound(request, response)
			return
		}

		for _, user := range users {
			resp.AddUser(&user)
		}
	}

	resp.Send(response)
}
