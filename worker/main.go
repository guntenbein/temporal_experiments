package main

import (
	"log"
	"temporal_experiments"
	"temporal_experiments/clients"
	localcontext "temporal_experiments/context"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
		ContextPropagators: []workflow.ContextPropagator{
			localcontext.NewStringMapPropagator([]string{temporal_experiments.CorrelationID}),
		},
		Logger: clients.Logger{},
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, temporal_experiments.QueueName, worker.Options{})

	w.RegisterWorkflow(temporal_experiments.MoveProductsWorkflow)
	w.RegisterWorkflow(temporal_experiments.MovePackagesWorkflow)
	w.RegisterActivity(temporal_experiments.NewMovingUnitsService(
		clients.BlockingScope{},
		clients.Search{},
		clients.Units{},
		clients.Index{Destination: "Search Service"},
		clients.Index{Destination: "Core-api"},
		clients.InMemoryCache{Data: make(map[string]interface{})},
		clients.Notify{},
	))

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
