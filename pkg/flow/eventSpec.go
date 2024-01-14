package flow

type EventSpec struct {
	Type           string             `json:"type,omitempty" validate:"required"`
	MonitorEnabled bool               `json:"monitorEnabled,omitempty"`
	EventActions   []EventActionsSpec `json:"eventActions,omitempty"`
	Data           interface{}        `json:"data,omitempty"`
	TransactionID  string             `json:"transaction_id,omitempty"`
}

func (e *EventSpec) ToEvent() Event {
	var eventActions []EventAction

	for _, action := range e.EventActions {
		eventActions = append(eventActions, action.ToEventAction())
	}

	return Event{
		Type:           e.Type,
		MonitorEnabled: e.MonitorEnabled,
		EventActions:   eventActions,
	}
}
