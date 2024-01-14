package main

import (
	"flow-events-connector/internal/config"
	"flow-events-connector/internal/database"
	"flow-events-connector/internal/flow"
	"flow-events-connector/internal/types"
	"fmt"
	"net/http"
	"time"

	cTypes "github.com/openfaas/connector-sdk/types"
	"github.com/openfaas/faas-provider/auth"
	"github.com/rs/zerolog/log"
)

const (
	DatastoreMigrationVersion = 1
)

func main() {
	cfg := config.NewConfig()
	creds := cTypes.GetCredentials()

	// Init DB
	db, err := database.InitDB(cfg, DatastoreMigrationVersion)
	if err != nil {
		log.Fatal().Msg("Failed to initialize database")
	}

	// Init Flow Service
	flowSvc := flow.NewService(db)

	log.Info().Msgf("Gateway URL: %s", cfg.Controller.GatewayURL)
	log.Info().Msgf("Rebuild interval: %s\tRebuild timeout: %s", cfg.Controller.RebuildInterval, cfg.Controller.RebuildTimeout)

	httpClient := cTypes.MakeClient(cfg.Controller.UpstreamTimeout)

	invoker := cTypes.NewInvoker(
		cfg.Controller.GatewayURL,
		httpClient,
		cfg.Controller.ContentType,
		cfg.Controller.PrintResponse,
		cfg.Controller.PrintRequestBody,
		cfg.Controller.UserAgent)

	go func() {
		for {
			r := <-invoker.Responses
			if r.Error != nil {
				log.Error().Msgf("Error with: %s, %s", r.Function, r.Error)
			} else {
				duration := fmt.Sprintf("%.2fs", r.Duration.Seconds())
				if r.Duration < time.Second*1 {
					duration = fmt.Sprintf("%dms", r.Duration.Milliseconds())
				}
				log.Info().Msgf("Response: %s [%d] (%s)", r.Function, r.Status, duration)
			}
		}
	}()

	if err := startEventsProbe(cfg, httpClient, creds, invoker, flowSvc); err != nil {
		log.Error().Msgf("Error: %s\n", err.Error())
	}
}

func startEventsProbe(cfg config.FlowEventsConnectorConfig, httpClient *http.Client, creds *auth.BasicAuthCredentials, invoker *cTypes.Invoker, flowSvc *flow.FlowService) error {
	ticker := time.NewTicker(cfg.Controller.RebuildInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		var events types.EventFunctions
		err := types.GetFunctionEvents(cfg, httpClient, creds, &events)
		if err != nil {
			log.Fatal().Msg("Could Not Get Function Events")
		}

		err = types.GetCoreEvents(cfg, &events)
		if err != nil {
			log.Fatal().Msg("Could Not Get Function Events")
		}
	}
}
