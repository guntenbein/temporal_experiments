package main

import (
	"context"
	"log"
	"temporal_experiments"
	"temporal_experiments/clients"
	localcontext "temporal_experiments/context"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
		ContextPropagators: []workflow.ContextPropagator{
			localcontext.NewStringMapPropagator([]string{temporal_experiments.CorrelationID}),
		},
		Logger: clients.Logger{},
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: temporal_experiments.QueueName,
	}
	// todo complete passing context
	ctx := context.WithValue(context.Background(), temporal_experiments.CorrelationID, "yahoo!!!")
	we, err := c.ExecuteWorkflow(ctx, workflowOptions, temporal_experiments.MoveProductsWorkflow,
		"company_1", "channel_1", "listing_1", "correlation_id_1", "searck_key_1",
		[]temporal_experiments.Move{{Type: "percentage", Mode: "loose", Destination: "listing_2", Value: 50}})
	if err != nil {
		panic(err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Wait for workflow completion. This is rarely needed in real use cases
	// when workflows are potentially long running
	var result []string
	err = we.Get(ctx, &result)
	if err != nil {
		panic(err)
	}
	log.Println("Finished workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
