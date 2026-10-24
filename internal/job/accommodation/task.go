package accommodation

import "time"

const (
	CrawlTaskType    = "accommodation.crawl"
	CrawlTaskVersion = 1
)

type CrawlTask struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Version      int       `json:"version"`
	ProviderCode string    `json:"provider_code"`
	Country      string    `json:"country,omitempty"`
	City         string    `json:"city,omitempty"`
	RequestedAt  time.Time `json:"requested_at"`
	TraceID      string    `json:"trace_id"`
}
