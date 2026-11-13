package accommodation

import (
	"time"
	"travel-crawler/internal/job"
)

const (
	CrawlTaskType    = "accommodation.crawl"
	CrawlTaskVersion = 1
)

type CrawlTask struct {
	job.TaskMetadata
	ProviderCode string     `json:"provider_code"`
	Country      string     `json:"country,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}
