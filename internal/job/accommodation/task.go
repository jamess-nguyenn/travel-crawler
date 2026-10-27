package accommodation

import (
	"travel-crawler/internal/job"
)

const (
	CrawlTaskType    = "accommodation.crawl"
	CrawlTaskVersion = 1
)

type CrawlTask struct {
	job.TaskMetadata
	ProviderCode string `json:"provider_code"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
}
