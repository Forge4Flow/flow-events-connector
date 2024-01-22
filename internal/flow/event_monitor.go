package flow

import (
	"flow-events-connector/internal/types"
	"fmt"

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
	for eventType, _ := range functions {
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

	reconnect := func(height uint64) {
		fmt.Printf("Reconnecting at block %d\n", height)

		data, errChan, err = client.SubscribeEventsByBlockHeight(svc.ctx, height, filter)
		if err != nil {
			log.Error().Msg(err.Error())
		}
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
				reconnect(lastHeight + 1)
				continue
			}

			for _, event := range eventData.Events {
				for _, eventFunction := range functions[event.Type] {
					eventFunction.InvokeFunction(svc.invoker)
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
				reconnect(lastHeight + 1)
				continue
			}

			fmt.Printf("~~~ ERROR: %s ~~~\n", err.Error())
			reconnect(lastHeight + 1)
			continue
		}
	}
}
