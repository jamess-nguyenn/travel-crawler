package kafka

import "time"

type FailureMessage struct {
	TaskID          string     `json:"task_id,omitempty"`
	TaskType        string     `json:"task_type,omitempty"`
	OriginalPayload []byte     `json:"original_payload"`
	Attempt         int        `json:"attempt"`
	MaxAttempt      int        `json:"max_attempt"`
	ErrorClass      string     `json:"error_class"`
	ErrorMessage    string     `json:"error_message"`
	FirstFailedAt   time.Time  `json:"first_failed_at"`
	FailedAt        time.Time  `json:"failed_at"`
	NextAttemptAt   *time.Time `json:"next_attempt_at,omitempty"`
	TraceID         string     `json:"trace_id"`
}
