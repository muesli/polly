package proposals

import (
	"net/http"
	"strconv"
	"time"

	"github.com/muesli/polly/api/db"
	"github.com/muesli/polly/api/utils"

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
func (r *ProposalResource) Put(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
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

	if auth.(db.User).ID != 1 && proposal.Starts.Before(time.Now()) {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Can't update a proposal once it has started",
			"ProposalResource PUT"))
	}

	proposal.Title = pps.Proposal.Title
	proposal.Description = pps.Proposal.Description
	proposal.Activities = pps.Proposal.Activities
	proposal.Contact = pps.Proposal.Contact
	proposal.Recipient = pps.Proposal.Recipient
	proposal.Recipient2 = pps.Proposal.Recipient2
	proposal.Value = pps.Proposal.Value
	proposal.RealValue = pps.Proposal.RealValue
	// proposal.Starts = pps.Proposal.Starts
	// proposal.FinishedDate = pps.Proposal.FinishedDate

	if auth.(db.User).ID == 1 {
		if !proposal.Moderated && pps.Proposal.Moderated {
			proposalAuthor, uerr := context.(*db.PollyContext).GetUserByID(proposal.UserID)
			if uerr != nil {
				panic(uerr)
			}
			utils.SendProposalAccepted(&proposalAuthor, &proposal)
		}

		proposal.Moderated = pps.Proposal.Moderated
	}

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
