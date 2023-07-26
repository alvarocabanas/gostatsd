package stats

import (
	"context"

	"github.com/atlassian/gostatsd"
)

// HeartBeater periodically sends a gauge for heartbeat purposes
type HeartBeater struct {
	metricName string
	tags       gostatsd.Tags
}

// NewHeartBeater creates a new HeartBeater
func NewHeartBeater(metricName string, tags gostatsd.Tags) *HeartBeater {
	return &HeartBeater{
		metricName: metricName,
		tags:       tags,
	}
}

// Run will run a HeartBeater in the background until the supplied context is closed.
func (hb *HeartBeater) Run(ctx context.Context) {
	statser := FromContext(ctx).WithTags(hb.tags)
	flushed, unregister := statser.RegisterFlush()
	defer unregister()

	reporter := ReportFromContext(ctx, statser)

	for {
		select {
		case <-ctx.Done():
			return
		case <-flushed:
			hb.emit(reporter)
		}
	}
}

func (hb *HeartBeater) emit(reporter ReportFunc) {
	reporter.Report(hb.metricName, 1, nil)
}
