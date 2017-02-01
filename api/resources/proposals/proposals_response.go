package proposals

import (
	"time"

	"github.com/muesli/polly/api/db"
	"github.com/muesli/polly/api/utils"

	"github.com/muesli/smolder"
)

// ProposalResponse is the common response to 'proposal' requests
type ProposalResponse struct {
	smolder.Response

	Proposals []proposalInfoResponse `json:"proposals,omitempty"`
	proposals []db.Proposal
}

type proposalInfoResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Activities   string    `json:"activities"`
	User         int64     `json:"user"`
	Contact      string    `json:"contact"`
	Recipient    string    `json:"recipient"`
	Recipient2   string    `json:"recipient2"`
	Value        uint64    `json:"value"`
	GrantType    string    `json:"granttype"`
	URL          string    `json:"url"`
	Starts       time.Time `json:"starts"`
	Ends         time.Time `json:"ends"`
	FinishedDate time.Time `json:"finished_date"`
	Ended        bool      `json:"ended"`
	Accepted     bool      `json:"accepted"`
	Moderated    bool      `json:"moderated"`
	Votes        uint64    `json:"votes"`
}

// Init a new response
func (r *ProposalResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Proposals = []proposalInfoResponse{}
}

// AddProposal adds a proposal to the response
func (r *ProposalResponse) AddProposal(proposal *db.Proposal) {
	r.proposals = append(r.proposals, *proposal)
	r.Proposals = append(r.Proposals, prepareProposalResponse(r.Context, proposal))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *ProposalResponse) EmptyResponse() interface{} {
	if len(r.proposals) == 0 {
		var out struct {
			Proposals interface{} `json:"proposals"`
		}
		out.Proposals = []proposalInfoResponse{}
		return out
	}
	return nil
}

func prepareProposalResponse(context smolder.APIContext, proposal *db.Proposal) proposalInfoResponse {
	ctx := context.(*db.PollyContext)
	resp := proposalInfoResponse{
		ID:           proposal.ID,
		Title:        proposal.Title,
		Description:  proposal.Description,
		Activities:   proposal.Activities,
		User:         proposal.UserID,
		Contact:      proposal.Contact,
		Recipient:    proposal.Recipient,
		Recipient2:   proposal.Recipient2,
		Value:        proposal.Value,
		Starts:       proposal.Starts,
		Ends:         proposal.Ends(ctx),
		FinishedDate: proposal.FinishedDate,
		Ended:        proposal.Ended(ctx),
		Votes:        proposal.Votes,
		Accepted:     proposal.Accepted(ctx),
		Moderated:    proposal.Moderated,
		URL:          utils.BuildURL(ctx.Config.Web.BaseURL, *proposal),
	}

	if proposal.Value < uint64(ctx.Config.App.Proposals.SmallGrantValueThreshold) {
		resp.GrantType = "small"
	} else {
		resp.GrantType = "large"
	}

	return resp
}
