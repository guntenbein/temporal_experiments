package temporal_experiments

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func MovePackagesWorkflow(ctx workflow.Context, companyID, uploadChannelID, sourceGroupID, processID, searchKey string, moves []Move) (err error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		// todo complete heartbeats
		HeartbeatTimeout: time.Second * 2,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	queryResult := "MOVING PACKAGES STARTED"
	err = workflow.SetQueryHandler(ctx, "state", func(input []byte) (string, error) {
		return queryResult, nil
	})
	if err != nil {
		return err
	}
	defer func() {
		// todo process errors
		workflow.ExecuteActivity(ctx, "UnblockScope", companyID, uploadChannelID, processID, "default").Get(ctx, nil)
		queryResult = "MOVING PACKAGES FINISHED"
	}()

	err = workflow.ExecuteActivity(ctx, "BlockScope", companyID, uploadChannelID, processID, "default").Get(ctx, nil)
	if err != nil {
		return err
	}
	queryResult = "SCOPE BLOCKED"

	err = workflow.ExecuteActivity(ctx, "CacheProductIDs", companyID, uploadChannelID, searchKey, processID).Get(ctx, nil)
	if err != nil {
		return err
	}
	queryResult = "IDS CACHED"

	err = workflow.ExecuteActivity(ctx, "MovePackages", companyID, uploadChannelID, sourceGroupID, processID, searchKey, moves).Get(ctx, nil)
	if err != nil {
		return err
	}
	queryResult = "PRODUCTS MOVED"

	futureIndexSearch, settableIndexSearch := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := workflow.ExecuteActivity(ctx, "IndexSearch", companyID, uploadChannelID, processID).Get(ctx, nil)
		if err != nil {
			settableIndexSearch.SetError(err)
			return
		}
		// todo what is the queryResult here?
		//err = workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "SEARCH INDEX COMPLETED").Get(ctx, nil)
		//if err != nil {
		//	settableIndexSearch.SetError(err)
		//	return
		//}
		// need to do this in order to finish the branch
		settableIndexSearch.SetValue(nil)
	})

	futureSyncListings, settableSyncListings := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := workflow.ExecuteActivity(ctx, "IndexListings", companyID, uploadChannelID, processID).Get(ctx, nil)
		if err != nil {
			settableSyncListings.SetError(err)
			return
		}
		// todo what is the queryResult here?
		//err = workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "SYNC LISTINGS COMPLETED").Get(ctx, nil)
		//if err != nil {
		//	settableSyncListings.SetError(err)
		//	return
		//}
		// need to do this in order to finish the branch
		settableSyncListings.SetValue(nil)
	})

	// return error from one of the parallel branches
	err = futureIndexSearch.Get(ctx, nil)
	if err != nil {
		return err
	}
	err = futureSyncListings.Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
