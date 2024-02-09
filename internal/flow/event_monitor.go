package flow

import (
	"flow-events-connector/internal/types"

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

	header, err := client.GetLatestBlockHeader(svc.ctx, true)
	if err != nil {
		return err
	}

	data, errChan, err := client.SubscribeEventsByBlockID(svc.ctx, header.ID, filter)
	if err != nil {
		return err
	}

	svc.WaitGroup.Add(1)

	reconnect := func(height uint64) error {
		log.Info().Msgf("Reconnecting at block %d\n", height)

		data, errChan, err = client.SubscribeEventsByBlockHeight(svc.ctx, height, filter)
		if err != nil {
			return err
		}

		return nil
	}

	lastHeight := header.Height
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
