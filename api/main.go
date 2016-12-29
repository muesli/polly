package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/muesli/smolder"
)

var (
	shutdownGracefully = false
	requestIncChan     chan int
)

func handleSignals() {
	pendingRequests := 0
	sigChan := make(chan os.Signal, 1)
	requestIncChan = make(chan int)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	boldGreen := string(byte(27)) + "[1;32m"
	boldRed := string(byte(27)) + "[1;31m"
	boldEnd := string(byte(27)) + "[0m"

	go func() {
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
}

func main() {
	handleSignals()

	parseSettings()
	SetupPostgres(config.Connections.PostgreSQLConnection)
	setupEmailTemplates()

	context := &PollyContext{}

	// Setup web-service
	smolderConfig := smolder.APIConfig{
		BaseURL:         config.API.BaseURL,
		PathPrefix:      config.API.PathPrefix,
		SwaggerAPIPath:  config.API.SwaggerAPIPath,
		SwaggerPath:     config.API.SwaggerPath,
		SwaggerFilePath: config.API.SwaggerFilePath,
	}
	wsContainer := smolder.NewSmolderContainer(smolderConfig, &shutdownGracefully, requestIncChan)
	func(resources ...smolder.APIResource) {
		for _, r := range resources {
			r.Register(wsContainer, smolderConfig, context)
		}
	}(
		&ProposalResource{},
		&UserResource{},
		&SessionResource{},
	)

	// GlobalLog("Starting polly web-api...")
	server := &http.Server{Addr: config.API.Bind, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
