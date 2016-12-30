package main

import (
	"net/http"
	"strconv"

	_ "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProposalPutStruct holds all values of an incoming PUT request
type ProposalPutStruct struct {
	ProposalPostStruct
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
func (r *ProposalResource) Put(context smolder.APIContext, request *restful.Request, response *restful.Response, auth interface{}) {
	if auth == nil || auth.(DbUser).ID != 1 {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"ProposalResource PUT"))
		return
	}

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

	proposal, err := context.(*PollyContext).GetProposalByID(int64(id))
	if err != nil {
		r.NotFound(request, response)
		return
	}

	proposal.Moderated = pps.Proposal.Moderated
	err = proposal.Update(context.(*PollyContext))
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
