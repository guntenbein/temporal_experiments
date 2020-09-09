package temporal_experiments

const (
	QueueName      = "moving_units"
	CorrelationID  = "correlationID"
	QueryTypeState = "state"
)

type BusinessError struct{}

func (ise BusinessError) Error() string {
	return "business non-retryable error"
}
