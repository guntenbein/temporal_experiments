package clients

import (
	"context"
	"log"
)

type Search struct{}

func (s Search) GetProductIDs(ctx context.Context, companyID, uploadChannelID, searchKey string) ([]string, error) {
	// todo change the logger to the standard one
	log.Println("search product ids", "company", companyID, "channel", uploadChannelID, "searchKey", searchKey)
	return []string{"product1", "product2", "proiduct3"}, nil
}

func (s Search) GetPackageIDs(ctx context.Context, companyID, uploadChannelID, searchKey string) ([]string, error) {
	log.Println("search package ids", "company", companyID, "channel", uploadChannelID, "searchKey", searchKey)
	return []string{"package1", "package2", "package3"}, nil
}
