package clients

import (
	"context"
	"log"
)

type BlockingScope struct {
}

func (b BlockingScope) BlockScope(ctx context.Context, companyID, resourceID, processID, action string) error {
	// todo change the logger to the standard one
	log.Println("block scope", "company", companyID, "resource", resourceID, "process", processID, "action", action)
	return nil
}

func (b BlockingScope) UnblockScope(ctx context.Context, companyID, resourceID, processID, action string) error {
	log.Println("unblock scope", "company", companyID, "resource", resourceID, "process", processID, "action", action)
	return nil
}
