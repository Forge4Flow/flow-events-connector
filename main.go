package main

import (
	"flow-events-connector/internal/config"
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

	// // Init DB
	// db, err := database.InitDB(cfg, DatastoreMigrationVersion)
	// if err != nil {
	// 	log.Fatal().Msg("Failed to initialize database")
	// }

	log.Info().Msgf("Gateway URL: %s", cfg.GatewayURL)
	log.Info().Msgf("Rebuild interval: %s\tRebuild timeout: %s", cfg.RebuildInterval, cfg.RebuildTimeout)

	httpClient := cTypes.MakeClient(cfg.UpstreamTimeout)

	invoker := cTypes.NewInvoker(
		cfg.GatewayURL,
		httpClient,
		cfg.ContentType,
		cfg.PrintResponse,
		cfg.PrintRequestBody,
		cfg.UserAgent)

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

	// Init Flow Service
	// flowSvc := flow.NewService(cfg.Flow, db, invoker)
	flowSvc := flow.NewService(cfg.Flow, invoker)

	if err := startEventsProbe(cfg, httpClient, creds, flowSvc); err != nil {
		log.Error().Msgf("Error: %s\n", err.Error())
	}

	flowSvc.WaitGroup.Wait()
	flowSvc.StopEventMonitors()
}

func startEventsProbe(cfg config.FlowEventsConnectorConfig, httpClient *http.Client, creds *auth.BasicAuthCredentials, flowSvc *flow.FlowService) error {
	ticker := time.NewTicker(cfg.RebuildInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		flowSvc.StopEventMonitors()

		var events types.Networks
		err := types.GetFunctionEvents(cfg, httpClient, creds, &events)
		if err != nil {
			log.Fatal().Msg("Could Not Get Function Events")
		}

		err = types.GetCoreEvents(cfg, &events)
		if err != nil {
			log.Fatal().Msg("Could Not Get Function Events")
		}

		flowSvc.StartEventMonitors(events)
	}
}
