package worker

import "fmt"

type InvalidTaskError struct {
	Dimension string
	Reason    string
}

func (e *InvalidTaskError) Error() string {
	if e.Dimension == "" {
		return fmt.Sprintf("invalid task: %s", e.Reason)
	}

	return fmt.Sprintf("invalid %s task: %s", e.Dimension, e.Reason)
}
