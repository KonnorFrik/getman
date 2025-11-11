package formatter

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/types"
)

func TestUnitFormatRequest_Simple(t *testing.T) {
	req := &types.Request{
		Method:  http.MethodGet,
		URL:     "http://example.com",
		Headers: make(map[string]string),
	}

	formatted := FormatRequest(req)
	if !strings.Contains(formatted, "GET") {
		t.Error("expected formatted request to contain method")
	}
	if !strings.Contains(formatted, "http://example.com") {
		t.Error("expected formatted request to contain URL")
	}
}

func TestUnitFormatRequest_WithHeaders(t *testing.T) {
	req := &types.Request{
		Method: http.MethodGet,
		URL:    "http://example.com",
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	formatted := FormatRequest(req)
	if !strings.Contains(formatted, "Accept") {
		t.Error("expected formatted request to contain headers")
	}
}

func TestUnitFormatRequest_WithAuth(t *testing.T) {
	req := &types.Request{
		Method:  http.MethodGet,
		URL:     "http://example.com",
		Headers: make(map[string]string),
		Auth: &types.Auth{
			Type:  "bearer",
			Token: "testtoken123",
		},
	}

	formatted := FormatRequest(req)
	if !strings.Contains(formatted, "bearer") {
		t.Error("expected formatted request to contain auth type")
	}
}

func TestUnitFormatRequest_WithBody(t *testing.T) {
	req := &types.Request{
		Method:  http.MethodPost,
		URL:     "http://example.com",
		Headers: make(map[string]string),
		Body: &types.RequestBody{
			Type:        "json",
			Content:     []byte(`{"key": "value"}`),
			ContentType: "application/json",
		},
	}

	formatted := FormatRequest(req)
	if !strings.Contains(formatted, "key") {
		t.Error("expected formatted request to contain body")
	}
}

func TestUnitPrintRequest(t *testing.T) {
	req := &types.Request{
		Method:  http.MethodGet,
		URL:     "http://example.com",
		Headers: make(map[string]string),
	}

	PrintRequest(req)
}

func TestUnitFormatExecutionResult(t *testing.T) {
	result := &types.ExecutionResult{
		CollectionName: "Test Collection",
		Environment:    "test",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalDuration:  time.Second,
		Requests: []*types.RequestExecution{
			{
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    "http://example.com",
				},
				Response: &types.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
				},
				Duration:  time.Millisecond * 100,
				Timestamp: time.Now(),
			},
		},
		Statistics: &types.Statistics{
			Total:   1,
			Success: 1,
			Failed:  0,
			AvgTime: time.Millisecond * 100,
			MinTime: time.Millisecond * 100,
			MaxTime: time.Millisecond * 100,
		},
	}

	formatted := FormatExecutionResult(result)
	if !strings.Contains(formatted, "Test Collection") {
		t.Error("expected formatted result to contain collection name")
	}
	if !strings.Contains(formatted, "test") {
		t.Error("expected formatted result to contain environment")
	}
}

func TestUnitFormatExecutionResult_WithErrors(t *testing.T) {
	result := &types.ExecutionResult{
		CollectionName: "Test Collection",
		Environment:    "test",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalDuration:  time.Second,
		Requests: []*types.RequestExecution{
			{
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    "http://example.com",
				},
				Error:     "request failed",
				Duration:  time.Millisecond * 100,
				Timestamp: time.Now(),
			},
		},
		Statistics: &types.Statistics{
			Total:   1,
			Success: 0,
			Failed:  1,
			AvgTime: time.Millisecond * 100,
			MinTime: time.Millisecond * 100,
			MaxTime: time.Millisecond * 100,
		},
	}

	formatted := FormatExecutionResult(result)
	if !strings.Contains(formatted, "request failed") {
		t.Error("expected formatted result to contain error")
	}
}

func TestUnitPrintExecutionResult(t *testing.T) {
	result := &types.ExecutionResult{
		CollectionName: "Test Collection",
		Environment:    "test",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalDuration:  time.Second,
		Requests:       []*types.RequestExecution{},
		Statistics: &types.Statistics{
			Total:   0,
			Success: 0,
			Failed:  0,
			AvgTime: 0,
			MinTime: 0,
			MaxTime: 0,
		},
	}

	PrintExecutionResult(result)
}

func TestUnitFormatStatistics(t *testing.T) {
	stats := &types.Statistics{
		Total:   10,
		Success: 8,
		Failed:  2,
		AvgTime: time.Millisecond * 100,
		MinTime: time.Millisecond * 50,
		MaxTime: time.Millisecond * 200,
	}

	formatted := FormatStatistics(stats)
	if !strings.Contains(formatted, "10") {
		t.Error("expected formatted statistics to contain total")
	}
	if !strings.Contains(formatted, "8") {
		t.Error("expected formatted statistics to contain success")
	}
	if !strings.Contains(formatted, "2") {
		t.Error("expected formatted statistics to contain failed")
	}
}

func TestUnitPrintStatistics(t *testing.T) {
	stats := &types.Statistics{
		Total:   10,
		Success: 8,
		Failed:  2,
		AvgTime: time.Millisecond * 100,
		MinTime: time.Millisecond * 50,
		MaxTime: time.Millisecond * 200,
	}

	PrintStatistics(stats)
}

func TestUnitMaskToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{"short token", "short", "***"},
		{"long token", "verylongtoken123456789", "very...6789"},
		{"exact 8 chars", "12345678", "1234...5678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskToken(tt.token)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

