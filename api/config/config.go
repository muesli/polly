package config

import (
	"errors"
	"flag"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Data holds all of polly's app settings
type Data struct {
	API struct {
		BaseURL         string
		PathPrefix      string
		Bind            string
		SwaggerAPIPath  string
		SwaggerPath     string
		SwaggerFilePath string
	}

	Connections struct {
		Email                EmailConfig
		PostgreSQLConnection PostgreSQLConnection
		Logging              struct {
			IrcHost       string
			IrcChannel    string
			IrcImpChannel string
		}
	}

	Web struct {
		BaseURL string
	}

	App struct {
		CryptPepper string

		Proposals struct {
			StartMonth                uint
			TotalRuntimeMonths        uint
			TotalGrantValue           uint
			GrantIntervalMonths       uint
			MaxGrantValue             uint
			MaxLargeGrantsPerMonth    uint
			SmallGrantValueThreshold  uint
			SmallGrantVetoThreshold   uint
			SmallGrantVoteThreshold   uint
			SmallGrantVoteRuntimeDays uint
		}

		Templates Templates
	}
}

// PostgreSQLConnection contains all of the db configuration values
type PostgreSQLConnection struct {
	User     string
	Password string
	Host     string
	DbName   string
	SslMode  string
}

// EmailConfig contains all email settings
type EmailConfig struct {
	AdminEmail string
	ReplyTo    string
	Mailman    struct {
		Name          string
		Address       string
		BounceAddress string
	}
	SMTP struct {
		User     string
		Password string
		Server   string
		Port     int
	}
	IMAP struct {
		User     string
		Password string
		Server   string
		Port     int
	}
}

// EmailTemplate holds all values of one email template
type EmailTemplate struct {
	Subject string
	Text    string
	HTML    string
}

// Templates holds all email templates
type Templates struct {
	Invitation         EmailTemplate
	ModerationProposal EmailTemplate
	ProposalAccepted   EmailTemplate
	ProposalStarted    EmailTemplate
}

var (
	// Settings contains the parsed configuration values
	Settings *Data
)

// ParseSettings parses the config file
func ParseSettings() {
	logLevelStr := flag.String("loglevel", "info", "Log level")
	configFile := flag.String("configfile", "config.json", "config file in the JSON format")
	flag.Parse()

	logLevel, err := log.ParseLevel(*logLevelStr)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	if configFile == nil || len(*configFile) == 0 {
		log.Panic(errors.New("Did not get a config file passed in"))
	}
	log.WithField("File", *configFile).Info("Using config file")

	// Parse config file
	configData := Data{}
	handler := NewHandler(*configFile, &configData, nil)
	if handler == nil {
		log.Fatal(errors.New("Config handler is nil, cannot continue"))
	}
	if !handler.LastReadValid() {
		log.WithField(
			"File",
			*configFile,
		).Fatal(errors.New("Did not get a valid config file"))
	}

	Settings = handler.CurrentData().(*Data)
	//FIXME catch empty conf fields
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
