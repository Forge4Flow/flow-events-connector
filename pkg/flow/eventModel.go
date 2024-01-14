package flow

import (
	"time"
)

type Model interface {
	GetID() int64
	GetType() string
	GetLastBlockHeight() uint64
	GetMonitorEnabled() bool
	GetEventActions() []EventAction
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() *time.Time
	ToEventSpec() *EventSpec
}

type Event struct {
	ID              int64  `mysql:"id" postgres:"id" sqlite:"id"`
	Type            string `mysql:"type" postgres:"type" sqlite:"type"`
	LastBlockHeight uint64 `mysql:"lastBlockHeight" postgres:"last_block_height" sqlite:"lastBlockHeight"`
	MonitorEnabled  bool   `mysql:"monitorEnabled" postgres:"monitor_enabled" sqlite:"monitorEnabled"`
	EventActions    []EventAction
	CreatedAt       time.Time  `mysql:"createdAt" postgres:"created_at" sqlite:"createdAt"`
	UpdatedAt       time.Time  `mysql:"updatedAt" postgres:"updated_at" sqlite:"updatedAt"`
	DeletedAt       *time.Time `mysql:"deletedAt" postgres:"deleted_at" sqlite:"deletedAt"`
}

func (event Event) GetID() int64 {
	return event.ID
}

func (event Event) GetType() string {
	return event.Type
}

func (event Event) GetLastBlockHeight() uint64 {
	return event.LastBlockHeight
}

func (event Event) GetMonitorEnabled() bool {
	return event.MonitorEnabled
}

func (event Event) GetEventActions() []EventAction {
	return event.EventActions
}

func (event Event) GetCreatedAt() time.Time {
	return event.CreatedAt
}

func (event Event) GetUpdatedAt() time.Time {
	return event.UpdatedAt
}

func (event Event) GetDeletedAt() *time.Time {
	return event.DeletedAt
}

func (event Event) ToEventSpec() *EventSpec {
	return &EventSpec{
		Type:           event.Type,
		MonitorEnabled: event.MonitorEnabled,
	}
}
