package reward

import "fmt"

// Error types for reward domain.
type RewardError struct {
	Code    string
	Message string
}

func (e *RewardError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrNotFound = &RewardError{Code: "REWARD_NOT_FOUND", Message: "reward does not exist"}
)
