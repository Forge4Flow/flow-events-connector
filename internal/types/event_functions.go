package types

import (
	"bytes"
	"encoding/json"
	"flow-events-connector/internal/config"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	cTypes "github.com/openfaas/connector-sdk/types"
	"github.com/openfaas/faas-provider/auth"
	"github.com/openfaas/faas-provider/sdk"
	ptypes "github.com/openfaas/faas-provider/types"
	"github.com/rs/zerolog/log"
)

type Networks map[string]EventFunctions

type EventFunctions map[string][]EventFunction

type EventFunction struct {
	FuncData  ptypes.FunctionStatus
	Name      string
	Namespace string
}

func (ef *EventFunction) String() string {
	if len(ef.Namespace) > 0 {
		return fmt.Sprintf("%s.%s", ef.Name, ef.Namespace)
	}

	return ef.Name
}

func (ef *EventFunction) InvokeFunction(cfg config.FlowEventsConnectorConfig, i *cTypes.Invoker, data interface{}) error {
	headers := http.Header{
		"X-Topic":     {cfg.Topic},
		"X-Connector": {cfg.ConnectorName},
	}

	// TODO: Allow Async Functions
	gwURL := fmt.Sprintf("%s/%s/%s", i.GatewayURL, "function", ef.String())

	// Serialize the data struct into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request body %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, gwURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create http request to %s %w", gwURL, err)
	}

	for k, v := range headers {
		req.Header[k] = v
	}

	if req.Body != nil {
		defer req.Body.Close()
	}
	start := time.Now()

	if cfg.PrintRequestBody {
		log.Info().Msgf("Request: %+v", req)
	}

	var body *[]byte
	res, err := i.Client.Do(req)
	if err != nil {
		i.Responses <- cTypes.InvokerResponse{
			Error:    fmt.Errorf("unable to invoke %s %w", ef.String(), err),
			Function: ef.Name,
			Topic:    "flow-events",
			Status:   http.StatusServiceUnavailable,
			Duration: time.Since(start),
		}
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
		bytesOut, err := io.ReadAll(res.Body)

		if err != nil {
			log.Printf("Error reading body")
			i.Responses <- cTypes.InvokerResponse{
				Error:    fmt.Errorf("unable to invoke %s %w", ef.String(), err),
				Status:   http.StatusServiceUnavailable,
				Function: ef.Name,
				Topic:    "flow-events",
				Duration: time.Since(start),
			}

			return fmt.Errorf("unable to read body %s", err)
		}

		body = &bytesOut
	}

	i.Responses <- cTypes.InvokerResponse{
		Body:     body,
		Status:   res.StatusCode,
		Header:   &res.Header,
		Function: ef.Name,
		Topic:    "flow-events",
		Duration: time.Since(start),
	}

	return nil
}

func GetFunctionEvents(c config.FlowEventsConnectorConfig, client *http.Client, creds *auth.BasicAuthCredentials, events *Networks) error {
	u, _ := url.Parse(c.GatewayURL)
	controller := sdk.NewSDK(u, creds, client)

	namespaces, err := controller.GetNamespaces()
	if err != nil {
		return err
	}

	if len(namespaces) == 0 {
		namespaces = []string{""}
	}

	for _, namespace := range namespaces {
		functions, err := controller.GetFunctions(namespace)
		if err != nil {
			return fmt.Errorf("unable to get functions in: %s, error: %w", namespace, err)
		}

		for _, function := range functions {
			err = toEventFunction(function, namespace, events, c.Topic)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetCoreEvents(c config.FlowEventsConnectorConfig, events *Networks) error {
	return nil
}

func toEventFunction(f ptypes.FunctionStatus, namespace string, events *Networks, topic string) error {
	if f.Annotations == nil {
		return fmt.Errorf("%s has no annotations", f.Name)
	}

	fTopic := (*f.Annotations)["topic"]
	fEvents := (*f.Annotations)["events"]

	if fTopic != topic {
		return fmt.Errorf("%s has wrong topic: %s", fTopic, f.Name)
	}

	if fEvents == "" {
		return fmt.Errorf("%s has no events defined", f.Name)
	}

	eventArray := strings.Split(fEvents, ",")
	for _, e := range eventArray {
		wNetwork := strings.SplitN(e, ".", 2)
		network := wNetwork[0]
		event := wNetwork[1]

		if _, ok := (*events)[network]; !ok {
			// If the network does not exist in the map, initialize its inner map
			(*events)[network] = make(EventFunctions)
		}

		if _, ok := (*events)[network][event]; !ok {
			// If the event does not exist in the network, initialize its slice
			(*events)[network][event] = []EventFunction{}
		}

		(*events)[network][event] = append((*events)[network][event], EventFunction{
			FuncData:  f,
			Name:      f.Name,
			Namespace: namespace,
		})
	}

	return nil
}
