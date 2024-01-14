package flow

import (
	"context"
	"errors"
	"fmt"

	user "github.com/forge4flow/forge4flow-core/pkg/authz/user"
	warrant "github.com/forge4flow/forge4flow-core/pkg/authz/warrant"
	"github.com/forge4flow/forge4flow-core/pkg/config"
	"github.com/forge4flow/forge4flow-core/pkg/service"

	"github.com/onflow/flow-go-sdk/access/http"
	"github.com/rs/zerolog/log"
)

type FlowService struct {
	service.BaseService
	Config       config.Forge4FlowConfig
	Repository   FlowEventRepository
	UserSvc      *user.UserService
	WarrantSvc   *warrant.WarrantService
	FlowClient   *http.Client
	queue        *Queue
	eventMonitor *EventMonitorService
}

func NewService(env service.Env, cfg config.Forge4FlowConfig, flowEventsRepo FlowEventRepository, userSvc *user.UserService, warrantSvc *warrant.WarrantService) *FlowService {
	flowClient, err := http.NewClient(cfg.GetFlowNetwork())
	if err != nil {
		log.Fatal().Err(err).Msg("Could not initialize and connect to the configured Flow Blockchain. Shutting down.")
	}

	svc := &FlowService{
		BaseService: service.NewBaseService(env),
		Config:      cfg,
		Repository:  flowEventsRepo,
		UserSvc:     userSvc,
		WarrantSvc:  warrantSvc,
		FlowClient:  flowClient,
	}

	svc.queue = newQueue(svc)
	go svc.queue.Start(25)

	svc.eventMonitor = newEventMonitorService(svc)
	events, err := svc.Repository.GetAllEvents(context.Background()) // Need to handle errors properly but it's probably safe to ignore them for now.
	if err != nil {
		fmt.Println(err)
	}

	for _, event := range events {
		fmt.Println("adding events")
		svc.eventMonitor.AddMonitor(event.GetType())
	}

	// eventChannel := svc.eventMonitor.eventChannel

	return svc
}

func (svc FlowService) ID() string {
	return service.FlowService
}

func (svc *FlowService) StartQueue() error {
	if !svc.queue.running {
		go svc.queue.Start(25)
		return nil
	}

	return errors.New("queue is already running")
}

func (svc *FlowService) StopQueue() error {
	if svc.queue.running {
		svc.queue.Stop()
		return nil
	}

	return errors.New("queue is not running")
}

func (svc *FlowService) Wait() {
	svc.queue.WaitGroup.Wait()
}

func (svc *FlowService) CreateQueueJob(job JobInterface) (string, error) {
	return svc.queue.CreateJob(job)
}

func (svc *FlowService) RemoveQueueJobByID(id string) error {
	return svc.queue.RemoveJobByID(id)
}

func resolver(str string) (string, error) {
	return str, nil
}
