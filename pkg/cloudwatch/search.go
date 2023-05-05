package cloudwatch

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/flanksource/apm-hub/api/logs"
)

func NewCloudWatchSearchBackend(config *logs.CloudWatchBackendConfig, client *cloudwatchlogs.Client) *cloudWatchSearch {
	return &cloudWatchSearch{
		client: client,
		config: config,
	}
}

type cloudWatchSearch struct {
	client *cloudwatchlogs.Client
	config *logs.CloudWatchBackendConfig
}

func (t *cloudWatchSearch) MatchRoute(q *logs.SearchParams) (match bool, isAdditive bool) {
	return t.config.CommonBackend.Routes.MatchRoute(q)
}

func (t *cloudWatchSearch) Search(q *logs.SearchParams) (logs.SearchResults, error) {
	logFilter := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: &t.config.LogGroup,
		Limit:        ptr(int32(q.Limit)),
		NextToken:    ptr(q.Page),
	}

	if q.GetStart() != nil {
		logFilter.StartTime = ptr(q.GetStart().UnixMilli())
	}

	if q.GetEnd() != nil {
		logFilter.EndTime = ptr(q.GetEnd().UnixMilli())
	}

	var result logs.SearchResults
	resp, err := t.client.FilterLogEvents(context.Background(), logFilter)
	if err != nil {
		return result, err
	}

	result.Results = make([]logs.Result, 0, len(resp.Events))
	for _, event := range resp.Events {
		result.Results = append(result.Results, logs.Result{
			Id:      deref(event.EventId),
			Message: deref(event.Message),
			Time:    time.UnixMilli(deref(event.Timestamp)).Format(time.RFC3339), // Convert from ms to time.RFC3339
			Labels:  t.config.Labels,
		})
	}

	result.NextPage = deref(resp.NextToken)

	return result, nil
}

func ptr[T any](val T) *T {
	return &val
}

func deref[T any](val *T) T {
	if val == nil {
		var v T
		return v
	}

	return *val
}
