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
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		Activities   string    `json:"activities"`
		Contact      string    `json:"contact"`
		Recipient    string    `json:"recipient"`
		Recipient2   string    `json:"recipient2"`
		Value        uint64    `json:"value"`
		RealValue    uint64    `json:"realvalue"`
		Moderated    bool      `json:"moderated"`
		Starts       time.Time `json:"starts"`
		FinishedDate time.Time `json:"finished_date"`
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
func (r *ProposalResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	authUser := db.User{}
	if auth, err := context.Authentication(request); err == nil {
		authUser = auth.(db.User)
	}

	resp := ProposalResponse{}
	resp.Init(context)

	pps := data.(*ProposalPostStruct)

	proposal := db.Proposal{
		UserID:       authUser.ID,
		Title:        pps.Proposal.Title,
		Description:  pps.Proposal.Description,
		Activities:   pps.Proposal.Activities,
		Contact:      pps.Proposal.Contact,
		Recipient:    pps.Proposal.Recipient,
		Recipient2:   pps.Proposal.Recipient2,
		Value:        pps.Proposal.Value,
		Starts:       pps.Proposal.Starts,
		FinishedDate: pps.Proposal.FinishedDate,
	}
	err := proposal.Save(context.(*db.PollyContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			err,
			"ProposalResource POST"))
		return
	}

	utils.SendModerationRequest(&proposal)

	resp.AddProposal(&proposal)
	resp.Send(response)
}
