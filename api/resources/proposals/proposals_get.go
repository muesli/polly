package proposals

import (
	"net/http"
	"strconv"

	"github.com/muesli/polly/api/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *ProposalResource) GetAuthRequired() bool {
	return true
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *ProposalResource) GetByIDsAuthRequired() bool {
	return true
}

// GetDoc returns the description of this API endpoint
func (r *ProposalResource) GetDoc() string {
	return "retrieve proposals"
}

// GetParams returns the parameters supported by this API endpoint
func (r *ProposalResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	// params = append(params, restful.QueryParameter("user_id", "id of a user").DataType("int64"))
	params = append(params, restful.QueryParameter("accepted", "only return accepted/rejected proposals").DataType("bool"))
	params = append(params, restful.QueryParameter("granttype", "small or large").DataType("string"))
	params = append(params, restful.QueryParameter("ended", "only returns finished proposals").DataType("bool"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *ProposalResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	authUser := db.User{}
	if auth, err := context.Authentication(request); err == nil {
		authUser = auth.(db.User)
	}

	resp := ProposalResponse{}
	resp.Init(context)

	ctx := context.(*db.PollyContext)
	for _, id := range ids {
		iid, err := strconv.Atoi(id)
		if err != nil {
			r.NotFound(request, response)
			return
		}
		proposal, err := ctx.GetProposalByID(int64(iid))
		if err != nil {
			r.NotFound(request, response)
			return
		}

		// only admin gets to see all proposals before moderation
		if authUser.ID == 1 || proposal.Moderated || proposal.Started(ctx) {
			resp.AddProposal(&proposal)
		} else {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				false,
				"Admin permission required for this operation",
				"ProposalResource GET"))
			return
		}
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *ProposalResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	authUser := db.User{}
	if auth, err := context.Authentication(request); err == nil {
		authUser = auth.(db.User)
	}

	accepted := params["accepted"]
	granttype := params["granttype"]
	ended := params["ended"]

	ctx := context.(*db.PollyContext)
	proposals, err := ctx.LoadAllProposals()
	if err != nil {
		r.NotFound(request, response)
		return
	}

	resp := ProposalResponse{}
	resp.Init(context)

	for _, proposal := range proposals {
		add := true
		// filter by grant-type
		if len(granttype) > 0 {
			if granttype[0] == "small" && proposal.Value >= uint64(ctx.Config.App.Proposals.SmallGrantValueThreshold) {
				add = false
			}
			if granttype[0] == "large" && proposal.Value < uint64(ctx.Config.App.Proposals.SmallGrantValueThreshold) {
				add = false
			}
		}
		if len(ended) > 0 {
			if ended[0] == "true" && !proposal.Ended(ctx) {
				// we only want proposals that ended already
				add = false
			}
		}
		if (len(ended) == 0 || ended[0] == "false") && proposal.Ended(ctx) {
			// we only want proposals that did not end yet
			add = false
		}
		if len(accepted) > 0 {
			if accepted[0] == "false" && proposal.Accepted(ctx) {
				add = false
			}
			if accepted[0] == "true" && !proposal.Accepted(ctx) {
				add = false
			}
		}

		// only admin gets to see all proposals before moderation
		if authUser.ID != 1 && (!proposal.Moderated || !proposal.Started(ctx)) {
			add = false
		}

		if add {
			resp.AddProposal(&proposal)
		}
	}

	resp.Send(response)
}
