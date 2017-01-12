package votes

import (
	"github.com/muesli/polly/api/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *VoteResource) GetAuthRequired() bool {
	return true
}

// GetDoc returns the description of this API endpoint
func (r *VoteResource) GetDoc() string {
	return "retrieve votes"
}

// GetParams returns the parameters supported by this API endpoint
func (r *VoteResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	// params = append(params, restful.QueryParameter("user_id", "id of a user").DataType("int64"))
	params = append(params, restful.QueryParameter("granttype", "small or large").DataType("string"))
	params = append(params, restful.QueryParameter("ended", "only returns finished votes").DataType("bool"))

	return params
}

// Get sends out items matching the query parameters
func (r *VoteResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	/*	authUser := db.User{}
		if auth, err := context.Authentication(request); err == nil {
			authUser = auth.(db.User)
		}*/

	ctx := context.(*db.PollyContext)
	votes, err := ctx.LoadAllUserVotes()
	if err != nil {
		r.NotFound(request, response)
		return
	}

	resp := VoteResponse{}
	resp.Init(context)

	for _, vote := range votes {
		resp.AddVote(&vote)
	}

	resp.Send(response)
}
