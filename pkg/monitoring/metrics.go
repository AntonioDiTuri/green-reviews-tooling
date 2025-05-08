// monitoring/benchmark.go
package monitoring

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/common/model"
)

// New struct to hold query results
type QueryResult struct {
	Query string
	Type  model.ValueType
	Value string
}

func ComputeBenchmarkingResults(ctx context.Context) ([]QueryResult, error) {
	q, err := NewQuery(
		WithClientTimeout(10 * time.Second),
	)
	if err != nil {
		return nil, err
	}

	queryMap := []struct {
		mType model.ValueType
		query string
		mVal  string
	}{
		// Not hardcode pod name
		{query: "rate(container_cpu_usage_seconds_total{pod='falco-driver-modern-ebpf-r9j4p'}[15m])"},
		{query: "avg_over_time(container_memory_rss{pod='falco-driver-modern-ebpf-r9j4p'}[15m])"},
		{query: "avg_over_time(container_memory_working_set_bytes{pod='falco-driver-modern-ebpf-r9j4p'}[15m])"},
	}

	results := make([]QueryResult, 0, len(queryMap))

	for idx := range queryMap {
		d, qErr := q.WithTimeRange(ctx, queryMap[idx].query, 15)
		if qErr != nil {
			return nil, qErr
		}

		queryMap[idx].mType = d.Type()
		queryMap[idx].mVal = d.String()

		// Append formatted result
		results = append(results, QueryResult{
			Query: queryMap[idx].query,
			Type:  queryMap[idx].mType,
			Value: strings.TrimSpace(queryMap[idx].mVal), // Clean up whitespace
		})
	}

	return results, nil
}
