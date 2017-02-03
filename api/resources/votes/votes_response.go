package votes

import (
	"github.com/muesli/polly/api/db"

	"github.com/muesli/smolder"
)

// VoteResponse is the common response to 'vote' requests
type VoteResponse struct {
	smolder.Response

	Votes []voteInfoResponse `json:"votes,omitempty"`
	votes []db.Vote
}

type voteInfoResponse struct {
	ID       int64 `json:"id"`
	User     int64 `json:"user"`
	Proposal int64 `json:"proposal"`
	Voted    bool  `json:"voted"`
}

// Init a new response
func (r *VoteResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Votes = []voteInfoResponse{}
}

// AddVote adds a vote to the response
func (r *VoteResponse) AddVote(vote db.Vote) {
	r.votes = append(r.votes, vote)
	r.Votes = append(r.Votes, prepareVoteResponse(r.Context, vote))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *VoteResponse) EmptyResponse() interface{} {
	if len(r.votes) == 0 {
		var out struct {
			Votes interface{} `json:"votes"`
		}
		out.Votes = []voteInfoResponse{}
		return out
	}
	return nil
}

func prepareVoteResponse(context smolder.APIContext, vote db.Vote) voteInfoResponse {
	resp := voteInfoResponse{
		ID:       vote.ID,
		User:     vote.UserID,
		Proposal: vote.ProposalID,
		Voted:    true,
	}

	return resp
}
