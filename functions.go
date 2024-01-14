package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/openfaas/connector-sdk/types"
	"github.com/openfaas/faas-provider/auth"
	"github.com/openfaas/faas-provider/sdk"
	ptypes "github.com/openfaas/faas-provider/types"
)

type EventFunctions map[string]Events

type Events map[string]EventFunction

type EventFunction struct {
	FuncData  ptypes.FunctionStatus
	Name      string
	Namespace string
}

func getFunctionEvents(c *types.ControllerConfig, client *http.Client, creds *auth.BasicAuthCredentials, events *EventFunctions) error {
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
			err = toEventFunction(function, namespace, events)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getCoreEvents(c *types.ControllerConfig, events *EventFunctions) error {
	return nil
}

func toEventFunction(f ptypes.FunctionStatus, namespace string, events *EventFunctions) error {
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

		(*events)[network][event] = EventFunction{
			FuncData:  f,
			Name:      f.Name,
			Namespace: namespace,
		}
	}

	return nil
}
