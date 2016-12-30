package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/muesli/polly/api/config"
	"github.com/muesli/polly/api/db"
	"github.com/muesli/polly/api/utils"

	"github.com/muesli/polly/api/resources/proposals"
	"github.com/muesli/polly/api/resources/sessions"
	"github.com/muesli/polly/api/resources/users"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful/swagger"
	"github.com/muesli/smolder"
)

func handleSignals() (chan int, bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	shutdownGracefully := false
	requestIncChan := make(chan int)

	go func() {
		boldGreen := string(byte(27)) + "[1;32m"
		boldRed := string(byte(27)) + "[1;31m"
		boldEnd := string(byte(27)) + "[0m"

		pendingRequests := 0
		for {
			select {
			case sig := <-sigChan:
				if !shutdownGracefully {
					shutdownGracefully = true
					fmt.Printf(boldGreen+"\nGot %s signal, shutting down gracefully. Press Ctrl-C again to stop now.\n\n"+boldEnd, sig.String())
					if pendingRequests == 0 {
						os.Exit(0)
					}
				} else {
					fmt.Printf(boldRed+"\nGot %s signal, shutting down now!\n\n"+boldEnd, sig.String())
					os.Exit(0)
				}

			case inc := <-requestIncChan:
				pendingRequests += inc
				if shutdownGracefully {
					log.Infoln("Pending requests:", pendingRequests)
					if pendingRequests == 0 {
						os.Exit(0)
					}
				}
			}
		}
	}()

	return requestIncChan, shutdownGracefully
}

func main() {
	ch, shutdownGracefully := handleSignals()

	config.ParseSettings()
	db.SetupPostgres(config.Settings.Connections.PostgreSQLConnection)
	utils.SetupEmailTemplates(*config.Settings)

	context := &db.PollyContext{
		Config: *config.Settings,
	}

	// Setup web-service
	smolderConfig := smolder.APIConfig{
		BaseURL:    config.Settings.API.BaseURL,
		PathPrefix: config.Settings.API.PathPrefix,
	}

	wsContainer := smolder.NewSmolderContainer(smolderConfig, &shutdownGracefully, ch)
	func(resources ...smolder.APIResource) {
		for _, r := range resources {
			r.Register(wsContainer, smolderConfig, context)
		}
	}(
		&sessions.SessionResource{},
		&users.UserResource{},
		&proposals.ProposalResource{},
	)

	if config.Settings.API.SwaggerFilePath != "" {
		wsConfig := swagger.Config{
			WebServices:     wsContainer.RegisteredWebServices(),
			WebServicesUrl:  config.Settings.API.BaseURL,
			ApiPath:         config.Settings.API.SwaggerAPIPath,
			SwaggerPath:     config.Settings.API.SwaggerPath,
			SwaggerFilePath: config.Settings.API.SwaggerFilePath,
		}
		swagger.RegisterSwaggerService(wsConfig, wsContainer)
	}

	// GlobalLog("Starting polly web-api...")
	server := &http.Server{Addr: config.Settings.API.Bind, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
