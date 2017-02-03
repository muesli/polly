package votes

import (
	"net/http"
	"strconv"

	"github.com/muesli/polly/api/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// VotePostStruct holds all values of an incoming POST request
type VotePostStruct struct {
	Vote struct {
		Proposal string `json:"proposal"`
		Voted    bool   `json:"voted"`
	} `json:"vote"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *VoteResource) PostAuthRequired() bool {
	return true
}

// PostDoc returns the description of this API endpoint
func (r *VoteResource) PostDoc() string {
	return "create a new vote"
}

// PostParams returns the parameters supported by this API endpoint
func (r *VoteResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *VoteResource) Post(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	ctx := context.(*db.PollyContext)
	authUser := db.User{}
	if auth, err := context.Authentication(request); err == nil {
		authUser = auth.(db.User)
	}

	resp := VoteResponse{}
	resp.Init(context)

	pps := VotePostStruct{}
	err := request.ReadEntity(&pps)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse POST data",
			"VoteResource POST"))
		return
	}

	pid, err := strconv.Atoi(pps.Vote.Proposal)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse proposal ID",
			"VoteResource POST"))
		return
	}

	proposal, err := ctx.GetProposalByID(int64(pid))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't find proposal",
			"VoteResource POST"))
	}

	vote, err := proposal.Vote(ctx, authUser)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't create vote",
			"VoteResource POST"))
	}

	resp.AddVote(vote)
	resp.Send(response)
}
