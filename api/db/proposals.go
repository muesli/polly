package db

import (
	"errors"
	"time"
)

// Proposal represents the db schema of a proposal
type Proposal struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Contact     string
	Recipient   string
	Recipient2  string
	Value       uint64
	Starts      time.Time
	Votes       uint64
	Moderated   bool
}

// LoadProposalByID loads a proposal by ID from the database
func (context *PollyContext) LoadProposalByID(id int64) (Proposal, error) {
	proposal := Proposal{}
	if id < 1 {
		return proposal, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, userid, title, description, contact, recipient, recipient2, value, starts, votes, moderated FROM proposals WHERE id = $1", id).Scan(&proposal.ID, &proposal.UserID, &proposal.Title, &proposal.Description, &proposal.Contact, &proposal.Recipient, &proposal.Recipient2, &proposal.Value, &proposal.Starts, &proposal.Votes, &proposal.Moderated)
	return proposal, err
}

// GetProposalByID returns a proposal by ID from the cache
func (context *PollyContext) GetProposalByID(id int64) (Proposal, error) {
	proposal := Proposal{}
	proposalCache, err := proposalsCache.Value(id, context)
	if err != nil {
		return proposal, err
	}

	proposal = *proposalCache.Data().(*Proposal)
	return proposal, nil
}

// LoadAllProposals loads all proposals from the database
func (context *PollyContext) LoadAllProposals() ([]Proposal, error) {
	proposals := []Proposal{}

	rows, err := context.Query("SELECT id, userid, title, description, contact, recipient, recipient2, value, starts, votes, moderated FROM proposals ORDER BY starts ASC")
	if err != nil {
		return proposals, err
	}

	defer rows.Close()
	for rows.Next() {
		proposal := Proposal{}
		err = rows.Scan(&proposal.ID, &proposal.UserID, &proposal.Title, &proposal.Description, &proposal.Contact, &proposal.Recipient, &proposal.Recipient2, &proposal.Value, &proposal.Starts, &proposal.Votes, &proposal.Moderated)
		if err != nil {
			return proposals, err
		}

		proposals = append(proposals, proposal)
	}

	return proposals, err
}

// Update a proposal in the database
func (proposal *Proposal) Update(context *PollyContext) error {
	_, err := context.Exec("UPDATE proposals SET title = $1, description = $2, contact = $3, recipient = $4, recipient2 = $5, value = $6, starts = $7, moderated = $8 WHERE id = $8", proposal.Title, proposal.Description, proposal.Contact, proposal.Recipient, proposal.Recipient2, proposal.Value, proposal.Starts, proposal.Moderated, proposal.ID)
	if err != nil {
		panic(err)
	}

	proposalsCache.Delete(proposal.ID)
	return err
}

// Save a proposal to the database
func (proposal *Proposal) Save(context *PollyContext) error {
	if proposal.Value > uint64(context.Config.App.Proposals.MaxGrantValue) {
		return errors.New("Grant value is too high")
	}

	if proposal.Value < uint64(context.Config.App.Proposals.SmallGrantValueThreshold) {
		if proposal.Value > uint64(context.SmallGrantMaxValue(uint(proposal.Ends(context).Month()))) {
			return errors.New("Proposal value is too high")
		}
	}

	if proposal.Starts.Before(time.Now()) {
		return errors.New("Invalid start date")
	}

	err := context.QueryRow("INSERT INTO proposals (userid, title, description, contact, recipient, recipient2, value, starts) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", proposal.UserID, proposal.Title, proposal.Description, proposal.Contact, proposal.Recipient, proposal.Recipient2, proposal.Value, proposal.Starts).Scan(&proposal.ID)
	proposalsCache.Delete(proposal.ID)
	return err
}

// Started returns true if a proposal has started
func (proposal *Proposal) Started(context *PollyContext) bool {
	return proposal.Starts.Before(time.Now())
}

// Ends returns when this proposal ends
func (proposal *Proposal) Ends(context *PollyContext) time.Time {
	return proposal.Starts.AddDate(0, 0, int(context.Config.App.Proposals.SmallGrantVoteRuntimeDays))
}

// Ended returns true if a proposal either ended or got rejected by votes
func (proposal *Proposal) Ended(context *PollyContext) bool {
	return proposal.Ends(context).Before(time.Now()) ||
		(proposal.Value < uint64(context.Config.App.Proposals.SmallGrantValueThreshold) &&
			proposal.Votes >= uint64(context.Config.App.Proposals.SmallGrantVoteThreshold))
}

// Accepted returns true if a proposal has finished and was accepted by poll
func (proposal *Proposal) Accepted(context *PollyContext) bool {
	return proposal.Ended(context) &&
		(proposal.Value >= uint64(context.Config.App.Proposals.SmallGrantValueThreshold) ||
			(proposal.Value < uint64(context.Config.App.Proposals.SmallGrantValueThreshold) &&
				proposal.Votes < uint64(context.Config.App.Proposals.SmallGrantVoteThreshold)))
}

// Vote marks a vote for a proposal
func (proposal *Proposal) Vote(context *PollyContext, user User) (Vote, error) {
	vote := Vote{
		UserID:     user.ID,
		ProposalID: proposal.ID,
		Vote:       true,
	}
	err := vote.Save(context)
	if err != nil {
		return Vote{}, err
	}

	err = context.QueryRow("UPDATE proposals SET votes=votes+1 WHERE id = $1 RETURNING votes", proposal.ID).Scan(&proposal.Votes)
	proposalsCache.Delete(proposal.ID)
	return vote, err
}
