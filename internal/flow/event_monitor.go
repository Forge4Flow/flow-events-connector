package flow

import (
	"flow-events-connector/internal/types"
	"fmt"
	"os"
	"strconv"

	flowSDK "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/rs/zerolog/log"
)

func (svc *FlowService) startEventMonitor(network string, functions types.EventFunctions) error {
	var client *grpc.Client
	switch network {
	case "mainnet":
		client = svc.mainnetClient
	case "testnet":
		client = svc.testnetClient
	case "crescendo":
		client = svc.crescendoClient
	case "emulator":
		client = svc.emulatorClient
	}

	// Create Event Filter
	var filter flowSDK.EventFilter
	for eventType := range functions {
		filter.EventTypes = append(filter.EventTypes, eventType)
	}

	// Read environment variable
	lastHeightString := os.Getenv("LAST_BLOCK_HEIGHT")
	var data <-chan flowSDK.BlockEvents
	var errChan <-chan error
	var err error
	var lastHeight uint64

	reconnect := func(height uint64) error {
		log.Info().Msgf("Reconnecting at block %d\n", height)

		data, errChan, err = client.SubscribeEventsByBlockHeight(svc.ctx, height, filter)
		if err != nil {
			return err
		}

		return nil
	}

	// Check if environment variable is set
	if lastHeightString == "" {
		header, err := client.GetLatestBlockHeader(svc.ctx, true)
		if err != nil {
			return err
		}

		data, errChan, err = client.SubscribeEventsByBlockID(svc.ctx, header.ID, filter)
		if err != nil {
			return err
		}

		lastHeight = header.Height
	} else {
		// Convert string to uint64
		lastHeight, err = strconv.ParseUint(lastHeightString, 10, 64)
		if err != nil {
			log.Error().Msgf("Error converting value to uint64: %v", err)
		}

		reconnect(lastHeight)
	}

	svc.WaitGroup.Add(1)

	for {
		select {
		case <-svc.ctx.Done():
			svc.WaitGroup.Done()
			return nil

		case eventData, ok := <-data:
			if !ok {
				if svc.ctx.Err() != nil {
					svc.WaitGroup.Done()
					return nil
				}
				// unexpected close
				err := reconnect(lastHeight + 1)
				if err != nil {
					svc.WaitGroup.Done()
					return err
				}

				continue
			}

			for _, event := range eventData.Events {
				for _, eventFunction := range functions[event.Type] {
					eventFunction.InvokeFunction(*svc.cfg, svc.invoker, event)
				}
			}

			err := os.Setenv("LAST_BLOCK_HEIGHT", fmt.Sprint(lastHeight))
			if err != nil {
				log.Error().Msgf("Error setting last height to environment variable: %v", err)
			}

			lastHeight = eventData.Height

		case err, ok := <-errChan:
			if !ok {
				if svc.ctx.Err() != nil {
					svc.WaitGroup.Done()
					return nil // graceful shutdown
				}
				// unexpected close
				err := reconnect(lastHeight + 1)
				if err != nil {
					svc.WaitGroup.Done()
					return err
				}

				continue
			}

			log.Info().Msgf("~~~ ERROR: %s ~~~\n", err.Error())
			err = reconnect(lastHeight + 1)
			if err != nil {
				svc.WaitGroup.Done()
				return err
			}

			continue
		}
	}
}
