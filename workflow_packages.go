package temporal_experiments

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func MovePackagesWorkflow(ctx workflow.Context, companyID, uploadChannelID, sourceGroupID, processID, searchKey string, moves []Move) (err error) {
	// todo use standard logger

	// todo separate local activity options should be provided
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		// todo complete heartbeats
		HeartbeatTimeout: time.Second * 2,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "MOVING PACKAGES STARTED").Get(ctx, nil)

	defer func() {
		// todo process errors
		workflow.ExecuteActivity(ctx, "UnblockScope", companyID, uploadChannelID, processID, "default").Get(ctx, nil)
		if err != nil {
			// todo process errors
			workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "MOVING PACKAGES FAILED").Get(ctx, nil)
		} else {
			workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "MOVING PACKAGES COMPLETED").Get(ctx, nil)
		}
	}()

	err = workflow.ExecuteActivity(ctx, "BlockScope", companyID, uploadChannelID, processID, "default").Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "SCOPE BLOCKED").Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "CacheProductIDs", companyID, uploadChannelID, searchKey, processID).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "IDS CACHED").Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "MovePackages", companyID, uploadChannelID, sourceGroupID, processID, searchKey, moves).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "PRODUCTS MOVED").Get(ctx, nil)
	if err != nil {
		return err
	}

	futureIndexSearch, settableIndexSearch := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := workflow.ExecuteActivity(ctx, "IndexSearch", companyID, uploadChannelID, processID).Get(ctx, nil)
		if err != nil {
			settableIndexSearch.SetError(err)
			return
		}
		err = workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "SEARCH INDEX COMPLETED").Get(ctx, nil)
		if err != nil {
			settableIndexSearch.SetError(err)
			return
		}
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
		err = workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "SYNC LISTINGS COMPLETED").Get(ctx, nil)
		if err != nil {
			settableSyncListings.SetError(err)
			return
		}
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
