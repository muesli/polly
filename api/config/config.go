package config

import (
	"errors"
	"flag"
	"reflect"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// ConfigData holds all of polly's app settings
type ConfigData struct {
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
			SmallGrantThreshold uint
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

// EmailTemplate holds all values of an email template
type EmailTemplate struct {
	Subject string
	Text    string
	HTML    string
}

type Templates struct {
	Invitation         EmailTemplate
	ModerationProposal EmailTemplate
}

type EmailConfig struct {
	AdminEmail string
	ReplyTo    string
	SMTP       struct {
		User     string
		Password string
		Server   string
		Port     int
	}
	IMAP struct {
		User        string
		Password    string
		Server      string
		Port        int
		LastMessage string
		LastSeen    int64
	}
}

var (
	configHandler *ConfigHandler
	Settings      *ConfigData
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
	configData := ConfigData{}
	if configHandler = NewConfigHandler(*configFile, &configData, nil); configHandler == nil {
		log.Fatal(errors.New("Config handler is nil, cannot continue"))
	}
	if !configHandler.LastReadValid() {
		log.WithField(
			"File",
			*configFile,
		).Fatal(errors.New("Did not get a valid config file"))
	}

	Settings = configHandler.CurrentData().(*ConfigData)
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
