package main

import (
	"flow-events-connector/internal/config"
	"fmt"
	"time"

	"github.com/openfaas/connector-sdk/types"
	"github.com/rs/zerolog/log"
)

const (
	userAgent = "forge4flow/flow-events-connector@v0.0.1"
	topic     = "flow-events"
)

func main() {
	cfg, err := config.GetControllerConfig()
	if err != nil {
		log.Fatal().Msg("Could Not Get Controller Config")
	}
	cfg.UserAgent = userAgent
	creds := types.GetCredentials()

	log.Info().Msgf("Gateway URL: %s", cfg.GatewayURL)
	log.Info().Msgf("Async Invocation: %v", cfg.AsyncFunctionInvocation)
	log.Info().Msgf("Rebuild interval: %s\tRebuild timeout: %s", cfg.RebuildInterval, cfg.RebuildInterval)

	httpClient := types.MakeClient(cfg.UpstreamTimeout)

	var events EventFunctions
	err = getFunctionEvents(cfg, httpClient, creds, &events)
	if err != nil {
		log.Fatal().Msg("Could Not Get Function Events")
	}

	err = getCoreEvents(cfg, &events)
	if err != nil {
		log.Fatal().Msg("Could Not Get Function Events")
	}

	invoker := types.NewInvoker(
		gatewayRoute(cfg),
		httpClient,
		cfg.ContentType,
		cfg.PrintResponse,
		cfg.PrintRequestBody,
		cfg.UserAgent)

	go func() {
		for {
			r := <-invoker.Responses
			if r.Error != nil {
				log.Error().Msgf("Error with: %s, %s", r.Function, err.Error())
			} else {
				duration := fmt.Sprintf("%.2fs", r.Duration.Seconds())
				if r.Duration < time.Second*1 {
					duration = fmt.Sprintf("%dms", r.Duration.Milliseconds())
				}
				log.Info().Msgf("Response: %s [%d] (%s)", r.Function, r.Status, duration)
			}
		}
	}()
}

func gatewayRoute(config *types.ControllerConfig) string {
	if config.AsyncFunctionInvocation {
		return fmt.Sprintf("%s/%s", config.GatewayURL, "async-function")
	}

	return fmt.Sprintf("%s/%s", config.GatewayURL, "function")
}
