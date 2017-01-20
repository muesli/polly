package proposals

import (
	"net/http"
	"time"

	"github.com/muesli/polly/api/db"
	"github.com/muesli/polly/api/utils"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProposalPostStruct holds all values of an incoming POST request
type ProposalPostStruct struct {
	Proposal struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Recipient   string    `json:"recipient"`
		Value       uint64    `json:"value"`
		Moderated   bool      `json:"moderated"`
		Starts      time.Time `json:"starts"`
	} `json:"proposal"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *ProposalResource) PostAuthRequired() bool {
	return true
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
func (r *ProposalResource) Post(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	authUser := db.User{}
	if auth, err := context.Authentication(request); err == nil {
		authUser = auth.(db.User)
	}

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
		UserID:      authUser.ID,
		Title:       pps.Proposal.Title,
		Description: pps.Proposal.Description,
		Recipient:   pps.Proposal.Recipient,
		Value:       pps.Proposal.Value,
		Ends:        pps.Proposal.Ends,
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
