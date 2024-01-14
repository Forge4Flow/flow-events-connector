package flow

import (
	"flow-events-connector/internal/database"
	"flow-events-connector/internal/service"
)

type FlowService struct {
	service.BaseService
}

func NewService(db database.Database) *FlowService {
	return &FlowService{
		BaseService: service.NewBaseService(db),
	}
}
