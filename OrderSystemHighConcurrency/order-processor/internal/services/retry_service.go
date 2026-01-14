package services

type RetryService struct {
	maxRetries int
}

func NewRetryService(maxRetries int) *RetryService {
	return &RetryService{
		maxRetries: maxRetries,
	}
}

func (r *RetryService) ShouldRetry(attempt int) bool {
	return attempt < r.maxRetries
}
