package db

import (
	"errors"
	"time"
)

// Proposal represents the db schema of a proposal
type Proposal struct {
	ID           int64
	UserID       int64
	Title        string
	Description  string
	Activities   string
	Contact      string
	Recipient    string
	Recipient2   string
	Value        uint64
	RealValue    uint64
	Starts       time.Time
	FinishedDate time.Time
	Votes        uint64
	Vetos        uint64
	Moderated    bool
	StartTrigger bool
}

// LoadProposalByID loads a proposal by ID from the database
func (context *PollyContext) LoadProposalByID(id int64) (Proposal, error) {
	proposal := Proposal{}
	if id < 1 {
		return proposal, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, userid, title, description, activities, contact, recipient, recipient2, value, realvalue, starts, votes, vetos, moderated, started, finisheddate FROM proposals WHERE id = $1", id).
		Scan(&proposal.ID, &proposal.UserID, &proposal.Title, &proposal.Description, &proposal.Activities, &proposal.Contact, &proposal.Recipient, &proposal.Recipient2, &proposal.Value, &proposal.RealValue, &proposal.Starts, &proposal.Votes, &proposal.Vetos, &proposal.Moderated, &proposal.StartTrigger, &proposal.FinishedDate)
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

	rows, err := context.Query("SELECT id, userid, title, description, activities, contact, recipient, recipient2, value, realvalue, starts, votes, vetos, moderated, started, finisheddate FROM proposals ORDER BY starts ASC")
	if err != nil {
		return proposals, err
	}

	defer rows.Close()
	for rows.Next() {
		proposal := Proposal{}
		err = rows.Scan(&proposal.ID, &proposal.UserID, &proposal.Title, &proposal.Description, &proposal.Activities, &proposal.Contact, &proposal.Recipient, &proposal.Recipient2, &proposal.Value, &proposal.RealValue, &proposal.Starts, &proposal.Votes, &proposal.Vetos, &proposal.Moderated, &proposal.StartTrigger, &proposal.FinishedDate)
		if err != nil {
			return proposals, err
		}

		proposals = append(proposals, proposal)
	}

	return proposals, err
}

// Update a proposal in the database
func (proposal *Proposal) Update(context *PollyContext) error {
	_, err := context.Exec("UPDATE proposals SET title = $1, description = $2, activities = $3, contact = $4, recipient = $5, recipient2 = $6, value = $7, realvalue = $8, starts = $9, moderated = $10, started = $11, finisheddate = $12 WHERE id = $13", proposal.Title, proposal.Description, proposal.Activities, proposal.Contact, proposal.Recipient, proposal.Recipient2, proposal.Value, proposal.RealValue, proposal.Starts, proposal.Moderated, proposal.StartTrigger, proposal.FinishedDate, proposal.ID)
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
			return errors.New("Proposal value is too high for this polling period")
		}

		if proposal.Starts.Before(time.Now()) {
			return errors.New("Invalid start date")
		}
	} else {
		largeGrantStartMonth := ((int(proposal.Starts.Month()) + int(context.Config.App.Proposals.StartMonth)) % int(context.Config.App.Proposals.GrantIntervalMonths)) + int(proposal.Starts.Month())
		startDate := time.Date(proposal.Starts.Year(), time.Month(largeGrantStartMonth), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
		proposal.Starts = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 12, 0, 0, 0, time.UTC)

		if proposal.Starts.Before(time.Now()) {
			return errors.New("Invalid start date")
		}
	}

	err := context.QueryRow("INSERT INTO proposals (userid, title, description, activities, contact, recipient, recipient2, value, realvalue, starts, finisheddate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id", proposal.UserID, proposal.Title, proposal.Description, proposal.Activities, proposal.Contact, proposal.Recipient, proposal.Recipient2, proposal.Value, proposal.Value, proposal.Starts, proposal.FinishedDate).Scan(&proposal.ID)
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

// Ended returns true if a proposal ended
func (proposal *Proposal) Ended(context *PollyContext) bool {
	return proposal.Ends(context).Before(time.Now())
}

func (proposal *Proposal) IsTopTwo(context *PollyContext) bool {
	rows, err := context.Query("SELECT id FROM proposals WHERE starts = $1 ORDER BY votes DESC LIMIT 2", proposal.Starts)
	if err != nil {
		return false
	}

	defer rows.Close()
	for rows.Next() {
		p := Proposal{}
		err = rows.Scan(&p.ID)
		if err != nil {
			return false
		}

		if proposal.ID == p.ID {
			return true
		}
	}

	return false
}

// Accepted returns true if a proposal has finished and was accepted by poll
func (proposal *Proposal) Accepted(context *PollyContext) bool {
	if !proposal.Ended(context) {
		return false
	}
	if proposal.Votes < uint64(context.Config.App.Proposals.SmallGrantVoteThreshold) {
		return false
	}

	if proposal.Value >= uint64(context.Config.App.Proposals.SmallGrantValueThreshold) {
		return proposal.IsTopTwo(context)
	} else {
		return proposal.Vetos < uint64(context.Config.App.Proposals.SmallGrantVetoThreshold)
	}
}

// Vote marks a vote for a proposal
func (proposal *Proposal) Vote(context *PollyContext, user User, up bool) (Vote, error) {
	vote := Vote{
		UserID:     user.ID,
		ProposalID: proposal.ID,
		Vote:       up,
	}
	err := vote.Save(context)
	if err != nil {
		return Vote{}, err
	}

	if up {
		err = context.QueryRow("UPDATE proposals SET votes=votes+1 WHERE id = $1 RETURNING votes", proposal.ID).Scan(&proposal.Votes)
	} else {
		err = context.QueryRow("UPDATE proposals SET vetos=vetos+1 WHERE id = $1 RETURNING vetos", proposal.ID).Scan(&proposal.Vetos)
	}
	proposalsCache.Delete(proposal.ID)
	return vote, err
}
