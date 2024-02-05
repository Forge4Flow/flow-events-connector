package flow

import (
	"context"
	"flow-events-connector/internal/config"
	"flow-events-connector/internal/service"
	"flow-events-connector/internal/types"
	"sync"

	"github.com/onflow/flow-go-sdk/access/grpc"
	cTypes "github.com/openfaas/connector-sdk/types"
	"github.com/rs/zerolog/log"
)

type FlowService struct {
	service.BaseService
	mainnetClient       *grpc.Client
	testnetClient       *grpc.Client
	crescendoClient     *grpc.Client
	emulatorClient      *grpc.Client
	invoker             *cTypes.Invoker
	ctx                 context.Context
	eventMonitorRunning bool

	WaitGroup *sync.WaitGroup
}

// func NewService(cfg *config.FlowConfig, db database.Database, invoker *cTypes.Invoker) *FlowService {
func NewService(cfg *config.FlowConfig, invoker *cTypes.Invoker) *FlowService {
	// Setup Mainnet Client
	mainnet, err := grpc.NewClient(cfg.MainnetAccessNode)
	if err != nil {
		log.Fatal().Msgf("Unable To Init Mainnet Flow Client: %s", err)
	}

	// Setup Testnet Client
	testnet, err := grpc.NewClient(cfg.TestnetAccessNode)
	if err != nil {
		log.Fatal().Msgf("Unable To Init Testnet Flow Client: %s", err)
	}

	// Setup Crescendo Client
	var crescendo *grpc.Client
	if cfg.UseCrescendo {
		crescendo, err = grpc.NewClient(cfg.CrescendoAccessNode)
		if err != nil {
			log.Fatal().Msgf("Unable To Init Crescendo Flow Client: %s", err)
		}
	}

	// Setup Emulator Client
	var emulator *grpc.Client
	if cfg.UseEmulator {
		emulator, err = grpc.NewClient(cfg.EmulatorAccessNode)
		if err != nil {
			log.Fatal().Msgf("Unable To Init Emulator Flow Client: %s", err)
		}
	}

	return &FlowService{
		// BaseService:     service.NewBaseService(db),
		BaseService:     service.NewBaseService(),
		mainnetClient:   mainnet,
		testnetClient:   testnet,
		crescendoClient: crescendo,
		emulatorClient:  emulator,
		invoker:         invoker,
		ctx:             context.Background(),

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
