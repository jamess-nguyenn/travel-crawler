package job

import "time"

type Metadata interface {
	GetID() string
	GetType() string
	GetVersion() int
	GetRequestedAt() time.Time
	GetTraceID() string
}

type TaskMetadata struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Version     int       `json:"version"`
	RequestedAt time.Time `json:"requested_at"`
	TraceID     string    `json:"trace_id"`
}

func (t TaskMetadata) GetID() string {
	return t.ID
}

func (t TaskMetadata) GetType() string {
	return t.Type
}

func (t TaskMetadata) GetVersion() int {
	return t.Version
}

func (t TaskMetadata) GetRequestedAt() time.Time {
	return t.RequestedAt
}

func (t TaskMetadata) GetTraceID() string {
	return t.TraceID
}
