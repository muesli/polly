package proposals

import (
	"net/http"

	"github.com/muesli/polly/api/db"
	"github.com/muesli/polly/api/utils"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProposalPostStruct holds all values of an incoming POST request
type ProposalPostStruct struct {
	Proposal struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Recipient   string `json:"recipient"`
		Value       uint64 `json:"value"`
		Moderated   bool   `json:"moderated"`
	} `json:"proposal"`
}

// PostDoc returns the description of this API endpoint
func (r *ProposalResource) PostDoc() string {
	return "create a new proposal"
}

// PostParams returns the parameters supported by this API endpoint
func (r *ProposalResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *ProposalResource) Post(context smolder.APIContext, request *restful.Request, response *restful.Response, auth interface{}) {
	resp := ProposalResponse{}
	resp.Init(context)

	pps := ProposalPostStruct{}
	err := request.ReadEntity(&pps)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse POST data",
			"ProposalResource POST"))
		return
	}

	proposal := db.Proposal{
		UserID:      auth.(db.User).ID,
		Title:       pps.Proposal.Title,
		Description: pps.Proposal.Description,
		Recipient:   pps.Proposal.Recipient,
		Value:       pps.Proposal.Value,
	}
	err = proposal.Save(context.(*db.PollyContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't create proposal",
			"ProposalResource POST"))
		return
	}

	utils.SendModerationRequest(&proposal)

	resp.AddProposal(&proposal)
	resp.Send(response)
}
