package temporal_experiments

import (
	"log"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MoveProductsWorkflow(ctx workflow.Context, companyID, uploadChannelID, sourceGroupID, processID, searchKey string, moves []Move) (err error) {
	value := ctx.Value(CorrelationID)
	log.Printf("Correlation ID in-workflow: %s", value)

	// todo separate local activity options should be provided
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		// todo complete heartbeats
		HeartbeatTimeout: time.Second * 2,
		// no retry (MaximumAttempts:    1)
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "MOVING PRODUCTS STARTED").Get(ctx, nil)

	defer func() {
		// todo process errors
		workflow.ExecuteActivity(ctx, "UnblockScope", companyID, uploadChannelID, processID, "default").Get(ctx, nil)
		if err != nil {
			// todo process errors
			workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "MOVING PRODUCTS FAILED").Get(ctx, nil)
		} else {
			workflow.ExecuteActivity(ctx, "NotifyState", companyID, processID, "MOVING PRODUCTS COMPLETED").Get(ctx, nil)
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

	moveProductsCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute * 1,
		// retry except InternalServerError from the root of the project
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Minute,
			MaximumAttempts:        5,
			NonRetryableErrorTypes: []string{"InternalServerError"},
		},
	})

	err = workflow.ExecuteActivity(moveProductsCtx, "MoveProducts", companyID, uploadChannelID, sourceGroupID, processID, searchKey, moves).Get(ctx, nil)
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
