package db

import (
	"time"
)

// Proposal represents the db schema of a proposal
type Proposal struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Recipient   string
	Value       uint64
	Ends        time.Time
	Votes       uint64
	Moderated   bool
}

// LoadProposalByID loads a proposal by ID from the database
func (context *PollyContext) LoadProposalByID(id int64) (Proposal, error) {
	proposal := Proposal{}
	if id < 1 {
		return proposal, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, userid, title, description, recipient, value, ends, votes, moderated FROM proposals WHERE id = $1", id).Scan(&proposal.ID, &proposal.UserID, &proposal.Title, &proposal.Description, &proposal.Recipient, &proposal.Value, &proposal.Ends, &proposal.Votes, &proposal.Moderated)
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

	rows, err := context.Query("SELECT id, userid, title, description, recipient, value, ends, votes, moderated FROM proposals ORDER BY ends ASC")
	if err != nil {
		return proposals, err
	}

	defer rows.Close()
	for rows.Next() {
		proposal := Proposal{}
		err = rows.Scan(&proposal.ID, &proposal.UserID, &proposal.Title, &proposal.Description, &proposal.Recipient, &proposal.Value, &proposal.Ends, &proposal.Votes, &proposal.Moderated)
		if err != nil {
			return proposals, err
		}

		proposals = append(proposals, proposal)
	}

	return proposals, err
}

// Update a proposal in the database
func (proposal *Proposal) Update(context *PollyContext) error {
	_, err := context.Exec("UPDATE proposals SET title = $1, description = $2, moderated = $3 WHERE id = $4", proposal.Title, proposal.Description, proposal.Moderated, proposal.ID)
	if err != nil {
		panic(err)
	}

	proposalsCache.Delete(proposal.ID)
	return err
}

// Save a proposal to the database
func (proposal *Proposal) Save(context *PollyContext) error {
	proposal.Ends = time.Now().AddDate(0, 0, 14)

	err := context.QueryRow("INSERT INTO proposals (userid, title, description, recipient, value, ends) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", proposal.UserID, proposal.Title, proposal.Description, proposal.Recipient, proposal.Value, proposal.Ends).Scan(&proposal.ID)
	proposalsCache.Delete(proposal.ID)
	return err
}
