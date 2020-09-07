package clients

import (
	"context"
	"log"
	"temporal_experiments"
)

type Index struct {
	Destination string
}

func (i Index) IndexScope(ctx context.Context, companyID, uploadChannelID string, scope temporal_experiments.DataUpdatedScope) error {
	// todo change the logger to the standard one
	log.Println("index "+i.Destination, "company", companyID, "channel", uploadChannelID, "scope", scope)
	return nil
}
