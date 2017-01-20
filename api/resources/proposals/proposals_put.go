package proposals

import (
	"net/http"
	"strconv"

	"github.com/muesli/polly/api/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProposalPutStruct holds all values of an incoming PUT request
type ProposalPutStruct struct {
	ProposalPostStruct
}

// PutAuthRequired returns true because all requests need authentication
func (r *ProposalResource) PutAuthRequired() bool {
	return true
}

// PutDoc returns the description of this API endpoint
func (r *ProposalResource) PutDoc() string {
	return "update an existing proposal"
}

// PutParams returns the parameters supported by this API endpoint
func (r *ProposalResource) PutParams() []*restful.Parameter {
	return nil
}

// Put processes an incoming PUT (update) request
func (r *ProposalResource) Put(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	resp := ProposalResponse{}
	resp.Init(context)

	pps := ProposalPutStruct{}
	err := request.ReadEntity(&pps)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse PUT data",
			"ProposalResource PUT"))
		return
	}

	id, err := strconv.Atoi(request.PathParameter("proposal-id"))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Invalid proposal id",
			"ProposalResource PUT"))
		return
	}

	proposal, err := context.(*db.PollyContext).GetProposalByID(int64(id))
	if err != nil {
		r.NotFound(request, response)
		return
	}

	auth, err := context.Authentication(request)
	if err != nil || (auth.(db.User).ID != 1 && auth.(db.User).ID != proposal.UserID) {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"ProposalResource PUT"))
		return
	}

	if auth.(db.User).ID == 1 {
		proposal.Moderated = pps.Proposal.Moderated
	}
	proposal.Title = pps.Proposal.Title
	proposal.Description = pps.Proposal.Description
	proposal.Recipient = pps.Proposal.Recipient
	proposal.Value = pps.Proposal.Value
	proposal.Starts = pps.Proposal.Starts

	err = proposal.Update(context.(*db.PollyContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't update proposal",
			"ProposalResource PUT"))
		return
	}

	resp.AddProposal(&proposal)
	resp.Send(response)
}
