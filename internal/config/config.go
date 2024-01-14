package config

import (
	"fmt"
	"os"
	"time"

	"github.com/openfaas/connector-sdk/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func GetControllerConfig() (*types.ControllerConfig, error) {
	// Configure logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.Level(zerolog.DebugLevel))
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	gURL, ok := os.LookupEnv("gateway_url")

	if !ok {
		return nil, fmt.Errorf("gateway_url environment variable not set")
	}

	asynchronousInvocation := true
	if val, exists := os.LookupEnv("asynchronous_invocation"); exists {
		asynchronousInvocation = (val == "1" || val == "true")
	}

	contentType := "text/plain"
	if v, exists := os.LookupEnv("content_type"); exists && len(v) > 0 {
		contentType = v
	}

	var printResponseBody bool
	if val, exists := os.LookupEnv("print_response_body"); exists {
		printResponseBody = (val == "1" || val == "true")
	}

	rebuildInterval := time.Second * 10

	if val, exists := os.LookupEnv("rebuild_interval"); exists {
		d, err := time.ParseDuration(val)
		if err != nil {
			return nil, err
		}
		rebuildInterval = d
	}

	return &types.ControllerConfig{
		RebuildInterval:         rebuildInterval,
		GatewayURL:              gURL,
		AsyncFunctionInvocation: asynchronousInvocation,
		ContentType:             contentType,
		PrintResponse:           true,
		PrintResponseBody:       printResponseBody,
		PrintRequestBody:        false,
	}, nil
}
