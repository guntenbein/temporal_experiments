package main

import (
	"context"
	"log"
	"temporal_experiments"

	"go.temporal.io/sdk/client"
)

func main() {

	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	resp, err := c.QueryWorkflow(context.Background(), temporal_experiments.CorrelationID, "", temporal_experiments.QueryTypeState)
	if err != nil {
		log.Fatalln("Unable to query workflow", err)
	}
	var result interface{}
	if err := resp.Get(&result); err != nil {
		log.Fatalln("Unable to decode query result", err)
	}
	log.Println("Received query result", "Result", result)
}
