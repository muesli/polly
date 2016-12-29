package main

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/muesli/cache2go"
	uuid "github.com/nu7hatch/gouuid"
)

var (
	pgDB   *sql.DB
	pgConn PostgreSQLConnection

	proposalsCache = cache2go.Cache("track")
	usersCache     = cache2go.Cache("user")

	// ErrInvalidID is the error returned when encountering an invalid database ID
	ErrInvalidID = errors.New("Invalid id")
)

// PostgreSQLConnection contains all of the db configuration values
type PostgreSQLConnection struct {
	User     string
	Password string
	Host     string
	DbName   string
	SslMode  string
}

// Marshal returns a "Connection String" with escaped values of all non-empty
// fields as described at http://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
func (c *PostgreSQLConnection) Marshal() string {
	val := reflect.ValueOf(c).Elem()
	var out string
	l := val.NumField()

	r := strings.NewReplacer(`'`, `\'`, `\`, `\\`)

	for i := 0; i < l; i++ {
		var fieldValue string

		switch f := val.Field(i).Interface().(type) {
		case string:
			fieldValue = f
		case int:
			fieldValue = strconv.Itoa(f)
		}
		fieldType := val.Type().Field(i).Name

		if len(fieldValue) > 0 {
			out += strings.ToLower(fieldType) + "='" + r.Replace(fieldValue) + "'"
			if i < l {
				out += " "
			}
		}
	}

	return out
}

// SetupPostgres sets the db configuration
func SetupPostgres(pc PostgreSQLConnection) {
	pgConn = pc
}

// GetDatabase connects to the database on first run and returns the existing
// connection on further calls
func GetDatabase() *sql.DB {
	if pgDB == nil {
		var err error
		pgDB, err = sql.Open("postgres", pgConn.Marshal())
		if err != nil {
			panic(err)
		}

		tables := []string{
			`CREATE TABLE IF NOT EXISTS users
				(
				  id          	bigserial 	PRIMARY KEY,
				  username    	text      	NOT NULL,
				  password		text		NOT NULL,
				  about       	text,
				  email       	text		NOT NULL,
				  activated   	bool		DEFAULT false,
				  authtoken   	text      	NOT NULL,
				  CONSTRAINT  	uk_username	UNIQUE (username),
				  CONSTRAINT  	uk_email 	UNIQUE (email)
				)`,
			`CREATE TABLE IF NOT EXISTS proposals
				(
				  id          	bigserial 	PRIMARY KEY,
				  userid      	bigserial 	NOT NULL,
				  title       	text      	NOT NULL,
				  description	text      	NOT NULL,
				  recipient		text		NOT NULL,
				  value			int			NOT NULL,
				  ends			timestamp	NOT NULL,
				  votes	      	int       	DEFAULT 0,
				  moderated     bool        DEFAULT false,
				  CONSTRAINT  	fk_user		FOREIGN KEY (userid) REFERENCES users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
				)`,
			`CREATE TABLE IF NOT EXISTS votes
				(
				  id          	bigserial			PRIMARY KEY,
				  userid    	bigserial			NOT NULL,
				  proposalid   	bigserial			NOT NULL,
				  vote			bool				NOT NULL,
				  CONSTRAINT  	uk_user_proposal	UNIQUE (userid, proposalid),
				  CONSTRAINT  	fk_user				FOREIGN KEY (userid) REFERENCES users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE,
				  CONSTRAINT  	fk_proposal			FOREIGN KEY (proposalid) REFERENCES proposals (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
				)`,
		}

		// FIXME: add IF NOT EXISTS to CREATE INDEX statements (coming in v9.5)
		// See: http://www.postgresql.org/docs/devel/static/sql-createindex.html
		indexes := []string{
			`CREATE INDEX idx_users_email ON users(email)`,
			`CREATE INDEX idx_proposals_moderated ON proposals(moderated)`,
			`CREATE INDEX idx_proposals_value ON proposals(value)`,
			`CREATE INDEX idx_proposals_userid ON proposals(userid)`,
			`CREATE INDEX idx_proposals_ends ON proposals(ends)`,
			`CREATE INDEX idx_votes_userid ON votes(userid)`,
			`CREATE INDEX idx_votes_proposalid ON votes(proposalid)`,
		}

		for _, v := range tables {
			fmt.Println("Creating table:", v)
			_, err = pgDB.Exec(v)
			if err != nil {
				panic(err)
			}
		}
		for _, v := range indexes {
			fmt.Println("Creating index:", v)
			_, err = pgDB.Exec(v)
			if err != nil && strings.Index(err.Error(), "already exists") < 0 {
				fmt.Println("Error:", err)
			}
		}

	}

	return pgDB
}

// WipeDatabase drops all database tables - use carefully!
func WipeDatabase() {
	// Commented out to prevent accidental usage

	/*
		drops := []string{
			`DROP TABLE votes`,
			`DROP TABLE proposals`,
			`DROP TABLE users`,
		}

		for _, v := range drops {
			fmt.Println("Dropping table:", v)
			_, err := pgDB.Exec(v)
			if err != nil {
				panic(err)
			}
		}
	*/
}

func init() {
	fmt.Println("db.init")
	initCaches()

	negativeInf := time.Time{}
	positiveInf, _ := time.Parse("2006", "3000")

	pq.EnableInfinityTs(negativeInf, positiveInf)
}

// UUID returns a new unique identifier
func UUID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	uuid := strings.Join(strings.Split(u.String(), "-"), "")
	return uuid, nil
}

func initCaches() {
	usersCache.SetAddedItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Now in users-cache:", item.Key().(string), item.Data().(*DbUser).Username)
	})
	usersCache.SetAboutToDeleteItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Deleting from users-cache:", item.Key().(string), item.Data().(*DbUser).Username, item.CreatedOn())
	})
	usersCache.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		if len(args) == 1 {
			if context, ok := args[0].(*PollyContext); ok {
				user, err := context.LoadUserByID(key.(int64))
				if err != nil {
					fmt.Println("usersCache ERROR for key", key, ":", err)
					return nil
				}

				entry := cache2go.CreateCacheItem(key, 10*time.Minute, &user)
				return &entry
			}
		}
		fmt.Println("Got no APIContext passed in")
		return nil
	})

	proposalsCache.SetAddedItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Now in proposals-cache:", item.Key().(string), item.Data().(*DbProposal).Title)
	})
	proposalsCache.SetAboutToDeleteItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Deleting from proposals-cache:", item.Key().(string), item.Data().(*DbProposal).Title, item.CreatedOn())
	})
	proposalsCache.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		if len(args) == 1 {
			if context, ok := args[0].(*PollyContext); ok {
				proposal, err := context.LoadProposalByID(key.(int64))
				if err != nil {
					fmt.Println("proposalsCache ERROR for key", key, ":", err)
					return nil
				}

				entry := cache2go.CreateCacheItem(key, 10*time.Minute, &proposal)
				return &entry
			}
		}
		fmt.Println("Got no APIContext passed in")
		return nil
	})
}
