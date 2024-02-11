package flow

import (
	"context"
	"flow-events-connector/internal/config"
	"flow-events-connector/internal/types"
	"sync"

	"github.com/onflow/flow-go-sdk/access/grpc"
	cTypes "github.com/openfaas/connector-sdk/types"
	"github.com/rs/zerolog/log"
)

type FlowService struct {
	mainnetClient       *grpc.Client
	testnetClient       *grpc.Client
	crescendoClient     *grpc.Client
	emulatorClient      *grpc.Client
	invoker             *cTypes.Invoker
	ctx                 context.Context
	cfg                 *config.FlowEventsConnectorConfig
	eventMonitorRunning bool

	WaitGroup *sync.WaitGroup
}

func NewService(cfg *config.FlowEventsConnectorConfig, invoker *cTypes.Invoker) *FlowService {
	// Setup Mainnet Client
	mainnet, err := grpc.NewClient(cfg.Flow.MainnetAccessNode)
	if err != nil {
		log.Fatal().Msgf("Unable To Init Mainnet Flow Client: %s", err)
	}

	// Setup Testnet Client
	testnet, err := grpc.NewClient(cfg.Flow.TestnetAccessNode)
	if err != nil {
		log.Fatal().Msgf("Unable To Init Testnet Flow Client: %s", err)
	}

	// Setup Crescendo Client
	var crescendo *grpc.Client
	if cfg.Flow.UseCrescendo {
		crescendo, err = grpc.NewClient(cfg.Flow.CrescendoAccessNode)
		if err != nil {
			log.Fatal().Msgf("Unable To Init Crescendo Flow Client: %s", err)
		}
	}

	// Setup Emulator Client
	var emulator *grpc.Client
	if cfg.Flow.UseEmulator {
		emulator, err = grpc.NewClient(cfg.Flow.EmulatorAccessNode)
		if err != nil {
			log.Fatal().Msgf("Unable To Init Emulator Flow Client: %s", err)
		}
	}

	return &FlowService{
		mainnetClient:   mainnet,
		testnetClient:   testnet,
		crescendoClient: crescendo,
		emulatorClient:  emulator,
		invoker:         invoker,
		ctx:             context.Background(),
		cfg:             cfg,

		WaitGroup: new(sync.WaitGroup),
	}
}

func (svc *FlowService) StartEventMonitors(monitors types.Networks) error {
	svc.eventMonitorRunning = true

	for network, functions := range monitors {
		err := svc.startEventMonitor(network, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *FlowService) StopEventMonitors() {
	if svc.eventMonitorRunning {
		svc.ctx.Done()
	}

	svc.eventMonitorRunning = false
}
