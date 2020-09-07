package temporal_experiments

const (
	QueueName = "moving_units"
)

type InternalServerError struct{}

func (ise InternalServerError) Error() string {
	return "unexpected internal non-retryable error"
}
