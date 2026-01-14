package contracts

type RetryService interface {
	ShouldRetry(key any) bool
}
