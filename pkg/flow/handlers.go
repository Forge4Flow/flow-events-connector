package flow

import (
	"net/http"

	"github.com/forge4flow/forge4flow-core/pkg/service"
	"github.com/gorilla/websocket"
)

func (svc FlowService) Routes() ([]service.Route, error) {
	return []service.Route{
		service.ForgeRoute{
			Pattern:                    "/v1/flow/events",
			Method:                     "GET",
			Handler:                    service.NewRouteHandler(svc, GetEventsHandler),
			OverrideAuthMiddlewareFunc: service.PassthroughAuthMiddleware,
		},

		service.ForgeRoute{
			Pattern:                    "/v1/flow/events",
			Method:                     "POST",
			Handler:                    service.NewRouteHandler(svc, AddEventMonitorHandler),
			OverrideAuthMiddlewareFunc: service.PassthroughAuthMiddleware,
		},

		service.ForgeRoute{
			Pattern:                    "/v1/flow/events",
			Method:                     "DELETE",
			Handler:                    service.NewRouteHandler(svc, RemoveEventMonitorHandler),
			OverrideAuthMiddlewareFunc: service.PassthroughAuthMiddleware,
		},
	}, nil
}

func AddEventMonitorHandler(svc FlowService, w http.ResponseWriter, r *http.Request) error {
	var event EventSpec
	err := service.ParseJSONBody(r.Body, &event)
	if err != nil {
		return err
	}

	if event.Type == "" {
		return service.NewMissingRequiredParameterError("Type")
	}

	err = svc.AddEventMonitor(r.Context(), event)
	if err != nil {
		return err
	}

	service.SendJSONResponse(w, event)

	return nil
}

func RemoveEventMonitorHandler(svc FlowService, w http.ResponseWriter, r *http.Request) error {
	var event EventSpec
	err := service.ParseJSONBody(r.Body, &event)
	if err != nil {
		return service.NewInvalidRequestError("Invalid JSON body")
	}

	if event.Type == "" {
		return service.NewMissingRequiredParameterError("Type")
	}

	return svc.RemoveEventMonitor(r.Context(), event)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

func isWebSocketRequest(r *http.Request) bool {
	connectionHeaders := r.Header.Get("Connection")
	upgradeHeaders := r.Header.Get("Upgrade")
	return connectionHeaders == "Upgrade" && upgradeHeaders == "websocket"
}

type SubscriptionRequest struct {
	EventTypes []string `json:"eventTypes,omitempty"`
	SendAll    bool     `json:"sendAll,omitempty"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func GetEventsHandler(svc FlowService, w http.ResponseWriter, r *http.Request) error {
	if isWebSocketRequest(r) {
		return GetEventsWSHandler(svc, w, r)
	} else {
		return GetEventsHTTPHandler(svc, w, r)
	}
}

func GetEventsHTTPHandler(svc FlowService, w http.ResponseWriter, r *http.Request) error {
	eventModels, err := svc.Repository.GetAllEvents(r.Context())
	if err != nil {
		return err
	}

	var events []EventSpec
	for _, event := range eventModels {
		events = append(events, *event.ToEventSpec())
	}

	service.SendJSONResponse(w, events)
	return nil
}

func GetEventsWSHandler(svc FlowService, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Register the event channel to receive events from the EventMonitorService
	eventChannel := svc.eventMonitor.eventChannel

	// Create a channel to send filtered events to the WebSocket client
	filteredEvents := make(chan EventSpec)

	var subscriptionReq SubscriptionRequest

	// Start a goroutine to filter events based on the client's subscription preferences
	go func() {
		for event := range eventChannel {
			// Check if the event type is in the client's subscribed event types
			if isSubscribed(event.Type, subscriptionReq.EventTypes) {
				filteredEvents <- event
			}

			if subscriptionReq.SendAll {
				filteredEvents <- event
			}
		}
	}()

	// Start a goroutine to send filtered events from the filteredEvents channel to the WebSocket client
	go func() {
		for event := range filteredEvents {
			err := conn.WriteJSON(event)
			if err != nil {
				// Handle error, e.g., log and break the loop
				break
			}
		}
	}()

	// The function will block here until the connection is closed
	for {
		// Read incoming messages from the WebSocket client
		err := conn.ReadJSON(&subscriptionReq)
		if err != nil {
			// Connection closed, so we break the loop
			break
		}

		if len(subscriptionReq.EventTypes) > 0 && subscriptionReq.SendAll {
			conn.WriteJSON(ErrorMessage{Error: "You cannot subscribe to all events and set a filter at the same time."})
			conn.Close()
		}

		if len(subscriptionReq.EventTypes) == 0 && !subscriptionReq.SendAll {
			conn.WriteJSON(ErrorMessage{Error: "You must either subscribe to all events or set a filter."})
			conn.Close()
		}
	}

	return nil
}

// Helper function to check if an event type is subscribed by the client
func isSubscribed(eventType string, subscribedEventTypes []string) bool {
	for _, subscribedType := range subscribedEventTypes {
		if subscribedType == eventType {
			return true
		}
	}
	return false
}
