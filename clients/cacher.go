package clients

import (
	"context"
	"fmt"
	"temporal_experiments"
)

const idsKeyTemplate = "units:%s:%s"
const scopeKeyTemplate = "scope:%s:%s"

type InMemoryCache struct {
	Data map[string]interface{}
}

func (i InMemoryCache) PutIDs(ctx context.Context, companyID, searchKey string, ids []string) error {
	i.Data[fmt.Sprintf(idsKeyTemplate, companyID, searchKey)] = ids
	return nil
}

func (i InMemoryCache) GetIDs(ctx context.Context, companyID, searchKey string) ([]string, error) {
	idsInterface, ok := i.Data[fmt.Sprintf(idsKeyTemplate, companyID, searchKey)]
	if !ok || idsInterface == nil {
		return nil, fmt.Errorf("error getting ids from cache")
	}
	return idsInterface.([]string), nil
}

func (i InMemoryCache) PutScope(ctx context.Context, companyID, processID string, scope temporal_experiments.DataUpdatedScope) error {
	i.Data[fmt.Sprintf(scopeKeyTemplate, companyID, processID)] = scope
	return nil
}

func (i InMemoryCache) GetScope(ctx context.Context, companyID, processID string) (temporal_experiments.DataUpdatedScope, error) {
	scopeInterface, ok := i.Data[fmt.Sprintf(scopeKeyTemplate, companyID, processID)]
	if !ok || scopeInterface == nil {
		return temporal_experiments.DataUpdatedScope{}, fmt.Errorf("error getting scope from cache")
	}
	return scopeInterface.(temporal_experiments.DataUpdatedScope), nil
}
