package streams

type streamError struct {
	message string
}

func (e *streamError) Error() string {
	return e.message
}

const (
	InvalidStreamIDError       = "invalid stream ID"
	XAddIDGreaterThanZeroError = "The ID specified in XADD must be greater than 0-0"
	XAddIDSmallerThanTopError  = "The ID specified in XADD is equal or smaller than the target stream top item"
)
