package main

import (
	"errors"
	"flag"

	log "github.com/Sirupsen/logrus"
)

// ConfigData holds all of polly's app settings
type ConfigData struct {
	API struct {
		BaseURL         string
		PathPrefix      string
		Bind            string
		CryptPepper     string
		SwaggerAPIPath  string
		SwaggerPath     string
		SwaggerFilePath string
	}

	Web struct {
		BaseURL string
	}

	Connections struct {
		Email struct {
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
		PostgreSQLConnection PostgreSQLConnection
	}

	Proposals struct {
		SmallGrantThreshold uint
	}

	Templates struct {
		Invitation struct {
			Subject string
			Text    string
			HTML    string
		}
		ModerationProposal struct {
			Subject string
			Text    string
			HTML    string
		}
	}

	Logging struct {
		IrcHost       string
		IrcChannel    string
		IrcImpChannel string
	}
}

var (
	configHandler *ConfigHandler
	config        *ConfigData
)

func parseSettings() {
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
	config = configHandler.CurrentData().(*ConfigData)
	//FIXME catch empty conf fields
}
