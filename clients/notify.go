package clients

import (
	"context"
	"log"
)

type Notify struct{}

func (n Notify) Notify(ctx context.Context, companyID, processID, status string) error {
	log.Println(">>>>>> NOTIFY", "company", companyID, "process", processID, "status", status)
	return nil
}
