package contracts

type RetryService interface {
	ShouldRetry(attempt int) bool
}
