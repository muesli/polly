package main

import (
	"net/http"
	"strconv"
	"time"

	_ "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetDoc returns the description of this API endpoint
func (r *ProposalResource) GetDoc() string {
	return "retrieve proposals"
}

// GetParams returns the parameters supported by this API endpoint
func (r *ProposalResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	// params = append(params, restful.QueryParameter("user_id", "id of a user").DataType("int64"))
	params = append(params, restful.QueryParameter("granttype", "small or large").DataType("string"))
	params = append(params, restful.QueryParameter("ended", "only returns finished proposals").DataType("bool"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *ProposalResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	auth, err := context.Authentication(request)
	if auth == nil || err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Invalid accesstoken",
			"ProposalResource GET"))
		return
	}
	authUser := auth.(DbUser)

	resp := ProposalResponse{}
	resp.Init(context)

	for _, id := range ids {
		iid, err := strconv.Atoi(id)
		if err != nil {
			r.NotFound(request, response)
			return
		}
		proposal, err := context.(*PollyContext).GetProposalByID(int64(iid))
		if err != nil {
			r.NotFound(request, response)
			return
		}

		// only admin gets to see all proposals before moderation
		if authUser.ID == 1 || proposal.Moderated {
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
	auth, err := context.Authentication(request)
	if auth == nil || err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Invalid accesstoken",
			"ProposalResource GET"))
		return
	}
	authUser := auth.(DbUser)

	granttype := params["granttype"]
	ended := params["ended"]

	proposals, err := context.(*PollyContext).LoadAllProposals()
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
			if granttype[0] == "small" && proposal.Value >= uint64(config.Proposals.SmallGrantThreshold) {
				add = false
			}
			if granttype[0] == "large" && proposal.Value < uint64(config.Proposals.SmallGrantThreshold) {
				add = false
			}
		}
		if len(ended) > 0 {
			if ended[0] == "true" && proposal.Ends.After(time.Now()) {
				add = false
			}
		}
		if (len(ended) == 0 || ended[0] == "false") && proposal.Ends.Before(time.Now()) {
			add = false
		}

		// only admin gets to see all proposals before moderation
		if authUser.ID != 1 && !proposal.Moderated {
			add = false
		}

		if add {
			resp.AddProposal(&proposal)
		}
	}

	resp.Send(response)
}
