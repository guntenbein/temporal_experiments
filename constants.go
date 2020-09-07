package temporal_experiments

const (
	QueueName     = "moving_units"
	CorrelationID = "correlationID"
)

type InternalServerError struct{}

func (ise InternalServerError) Error() string {
	return "unexpected internal non-retryable error"
}
