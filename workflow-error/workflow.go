package workflow_error

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type BusinessError struct{}

func (be BusinessError) Error() string {
	return "business error"
}

func ExperimentalWorkflow(ctx workflow.Context) error {
	ctx = defineActivityOptions(ctx, "some")
	err := workflow.ExecuteActivity(ctx, "Activity1").Get(ctx, nil)
	if err != nil {
		return err
	}
	err = workflow.ExecuteActivity(ctx, "Activity2").Get(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func defineActivityOptions(ctx workflow.Context, queue string) workflow.Context {
	ao := workflow.ActivityOptions{
		TaskQueue:              queue,
		ScheduleToStartTimeout: time.Second,
		StartToCloseTimeout:    24 * time.Hour,
		HeartbeatTimeout:       time.Second * 2,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Minute * 5,
			NonRetryableErrorTypes: []string{"BusinessError"},
		},
	}
	ctxOut := workflow.WithActivityOptions(ctx, ao)
	return ctxOut
}
