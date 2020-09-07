package clients

import (
	"context"
	"log"
	"temporal_experiments"
	"time"
)

type Units struct{}

func (u Units) MoveProducts(ctx context.Context, companyID, uploadChannelID, sourceGroupID string, moves []temporal_experiments.Move, productIDs []string) (temporal_experiments.DataUpdatedScope, error) {
	log.Println("move products", "company", companyID, "channel", uploadChannelID, "source", sourceGroupID, "moves", moves, "products", productIDs)
	// imitate processing
	time.Sleep(time.Second * 10)
	return temporal_experiments.DataUpdatedScope{
		GroupIDList:   nil,
		PackageIDList: nil,
		ProductIDList: productIDs,
	}, nil
}

func (u Units) MovePackages(ctx context.Context, companyID, uploadChannelID, sourceGroupID string, moves []temporal_experiments.Move, packageIDs []string) (temporal_experiments.DataUpdatedScope, error) {
	log.Println("move packages", "company", companyID, "channel", uploadChannelID, "source", sourceGroupID, "moves", moves, "packages", packageIDs)
	// imitate processing
	time.Sleep(time.Second * 10)
	return temporal_experiments.DataUpdatedScope{
		GroupIDList:   nil,
		PackageIDList: packageIDs,
		ProductIDList: nil,
	}, nil
}
