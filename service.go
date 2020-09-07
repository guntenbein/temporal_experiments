package temporal_experiments

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type ScopeBlocker interface {
	BlockScope(ctx context.Context, companyID, resourceID, processID, action string) error
	UnblockScope(ctx context.Context, companyID, resourceID, processID, action string) error
}

type SearchedIDsProvider interface {
	GetProductIDs(ctx context.Context, companyID, uploadChannelID, searchKey string) ([]string, error)
	GetPackageIDs(ctx context.Context, companyID, uploadChannelID, searchKey string) ([]string, error)
}

type UnitsManager interface {
	MoveProducts(ctx context.Context, companyID, uploadChannelID, sourceGroupID string, moves []Move, productIDs []string) (DataUpdatedScope, error)
	MovePackages(ctx context.Context, companyID, uploadChannelID, sourceGroupID string, moves []Move, packageIDs []string) (DataUpdatedScope, error)
}

type Indexer interface {
	IndexScope(ctx context.Context, companyID, uploadChannelID string, scope DataUpdatedScope) error
}

type Cacher interface {
	PutIDs(ctx context.Context, companyID, searchKey string, ids []string) error
	GetIDs(ctx context.Context, companyID, searchKey string) ([]string, error)
	PutScope(ctx context.Context, companyID, processID string, scope DataUpdatedScope) error
	GetScope(ctx context.Context, companyID, processID string) (DataUpdatedScope, error)
}

type Notifier interface {
	Notify(ctx context.Context, companyID, processID, status string) error
}

type MovingUnitsService struct {
	scopeBlocker  ScopeBlocker
	idsProvider   SearchedIDsProvider
	unitsManager  UnitsManager
	searchIndexer Indexer
	syncListings  Indexer
	caher         Cacher
	notifier      Notifier
}

func NewMovingUnitsService(scopeBlocker ScopeBlocker,
	idsProvider SearchedIDsProvider,
	unitsManager UnitsManager,
	searchIndexer Indexer,
	syncListings Indexer,
	caher Cacher,
	notifier Notifier) *MovingUnitsService {
	return &MovingUnitsService{
		scopeBlocker:  scopeBlocker,
		idsProvider:   idsProvider,
		unitsManager:  unitsManager,
		searchIndexer: searchIndexer,
		syncListings:  syncListings,
		caher:         caher,
		notifier:      notifier,
	}
}

func (s *MovingUnitsService) BlockScope(ctx context.Context, companyID, resourceID, processID, action string) error {
	return s.scopeBlocker.BlockScope(ctx, companyID, resourceID, processID, action)
}

func (s *MovingUnitsService) UnblockScope(ctx context.Context, companyID, resourceID, processID, action string) error {
	return s.scopeBlocker.UnblockScope(ctx, companyID, resourceID, processID, action)
}

func (s *MovingUnitsService) CacheProductIDs(ctx context.Context, companyID, uploadChannelID, searchKey, processID string) error {
	ids, err := s.idsProvider.GetProductIDs(ctx, companyID, uploadChannelID, searchKey)
	if err != nil {
		return err
	}
	return s.caher.PutIDs(ctx, companyID, searchKey, ids)
}

func (s *MovingUnitsService) CachePackageIDs(ctx context.Context, companyID, uploadChannelID, searchKey, processID string) error {
	ids, err := s.idsProvider.GetPackageIDs(ctx, companyID, uploadChannelID, searchKey)
	if err != nil {
		return err
	}
	return s.caher.PutIDs(ctx, companyID, searchKey, ids)
}

func (s *MovingUnitsService) MoveProducts(ctx context.Context, companyID, uploadChannelID, sourceGroupID, processID, searchKey string, moves []Move) error {
	cancelc := make(chan struct{})
	defer close(cancelc)
	go func() {
		for {
			select {
			case <-cancelc:
				return
			default:
				{
					time.Sleep(1 * time.Second)
					activity.RecordHeartbeat(ctx)
				}
			}
		}
	}()
	ids, err := s.caher.GetIDs(ctx, companyID, searchKey)
	if err != nil {
		return err
	}
	scope, err := s.unitsManager.MoveProducts(ctx, companyID, uploadChannelID, sourceGroupID, moves, ids)
	if err != nil {
		return err
	}
	return s.caher.PutScope(ctx, companyID, processID, scope)
}

func (s *MovingUnitsService) MovePackages(ctx context.Context, companyID, uploadChannelID, sourceGroupID, processID, searchKey string, moves []Move) error {
	cancelc := make(chan struct{})
	defer close(cancelc)
	go func() {
		for {
			select {
			case <-cancelc:
				return
			default:
				{
					time.Sleep(1 * time.Second)
					activity.RecordHeartbeat(ctx)
				}
			}
		}
	}()
	ids, err := s.caher.GetIDs(ctx, companyID, searchKey)
	if err != nil {
		return err
	}
	scope, err := s.unitsManager.MovePackages(ctx, companyID, uploadChannelID, sourceGroupID, moves, ids)
	if err != nil {
		return err
	}
	return s.caher.PutScope(ctx, companyID, processID, scope)
}

func (s *MovingUnitsService) IndexSearch(ctx context.Context, companyID, uploadChannelID, processID string) error {
	scope, err := s.caher.GetScope(ctx, companyID, processID)
	if err != nil {
		return err
	}
	return s.searchIndexer.IndexScope(ctx, companyID, uploadChannelID, scope)
}

func (s *MovingUnitsService) IndexListings(ctx context.Context, companyID, uploadChannelID, processID string) error {
	scope, err := s.caher.GetScope(ctx, companyID, processID)
	if err != nil {
		return err
	}
	return s.syncListings.IndexScope(ctx, companyID, uploadChannelID, scope)
}

func (s *MovingUnitsService) NotifyState(ctx context.Context, companyID, processID, status string) error {
	return s.notifier.Notify(ctx, companyID, processID, status)
}
