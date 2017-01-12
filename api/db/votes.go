package db

// Vote represents the db schema of a vote
type Vote struct {
	ID         int64
	UserID     int64
	ProposalID int64
	Vote       bool
}

// LoadAllUserVotes loads all votes for a user from the database
func (context *PollyContext) LoadAllUserVotes() ([]Vote, error) {
	votes := []Vote{}

	rows, err := context.Query("SELECT id, userid, proposalid, vote FROM votes")
	if err != nil {
		return votes, err
	}

	defer rows.Close()
	for rows.Next() {
		vote := Vote{}
		err = rows.Scan(&vote.ID, &vote.UserID, &vote.ProposalID, &vote.Vote)
		if err != nil {
			return votes, err
		}

		votes = append(votes, vote)
	}

	return votes, err
}

// Save a proposal to the database
func (vote *Vote) Save(context *PollyContext) error {
	err := context.QueryRow("INSERT INTO votes (userid, proposalid, vote) VALUES ($1, $2, $3) RETURNING id", vote.UserID, vote.ProposalID, vote.Vote).Scan(&vote.ID)
	return err
}
