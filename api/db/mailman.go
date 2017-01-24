package db

// Mailman represents the db schema of mailman info
type Mailman struct {
	Mailbox  string
	LastSeen uint64
}

// GetMailman returns the mailman info for a mailbox
func (context *PollyContext) GetMailman(mailbox string) (Mailman, error) {
	mm := Mailman{}
	err := context.QueryRow("SELECT mailbox, lastseen FROM mailman WHERE mailbox = $1", mailbox).Scan(&mm.Mailbox, &mm.LastSeen)
	mm.Mailbox = mailbox
	return mm, err
}

// Update mailman info in the database
func (mm *Mailman) Update(context *PollyContext) error {
	_, err := context.Exec("INSERT INTO mailman (mailbox, lastseen) VALUES ($1, $2) "+
		"ON CONFLICT (mailbox) DO UPDATE SET lastseen = $2", mm.Mailbox, mm.LastSeen)
	if err != nil {
		panic(err)
	}
	return err
}
