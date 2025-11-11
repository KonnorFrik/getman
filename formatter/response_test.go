package formatter

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/types"
)

func TestUnitFormatResponse_Simple(t *testing.T) {
	resp := &types.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Headers:    make(map[string][]string),
		Body:       []byte("OK"),
		Duration:   time.Millisecond * 100,
		Size:       2,
	}

	formatted := FormatResponse(resp)
	if !strings.Contains(formatted, "200 OK") {
		t.Error("expected formatted response to contain status code")
	}
	if !strings.Contains(formatted, "OK") {
		t.Error("expected formatted response to contain body")
	}
}

func TestUnitFormatResponse_WithHeaders(t *testing.T) {
	resp := &types.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"X-Custom":     {"value"},
		},
		Body:     []byte("{}"),
		Duration: time.Millisecond * 100,
		Size:     2,
	}

	formatted := FormatResponse(resp)
	if !strings.Contains(formatted, "Content-Type") {
		t.Error("expected formatted response to contain headers")
	}
}

func TestUnitFormatResponse_WithJSONBody(t *testing.T) {
	resp := &types.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Headers:    make(map[string][]string),
		Body:       []byte(`{"key": "value"}`),
		Duration:   time.Millisecond * 100,
		Size:       15,
	}

	formatted := FormatResponse(resp)
	if !strings.Contains(formatted, "key") {
		t.Error("expected formatted response to contain JSON body")
	}
}

func TestUnitFormatResponse_WithTextBody(t *testing.T) {
	resp := &types.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Headers:    make(map[string][]string),
		Body:       []byte("text body"),
		Duration:   time.Millisecond * 100,
		Size:       9,
	}

	formatted := FormatResponse(resp)
	if !strings.Contains(formatted, "text body") {
		t.Error("expected formatted response to contain text body")
	}
}

func TestUnitFormatResponse_EmptyBody(t *testing.T) {
	resp := &types.Response{
		StatusCode: http.StatusNoContent,
		Status:     "204 No Content",
		Headers:    make(map[string][]string),
		Body:       []byte{},
		Duration:   time.Millisecond * 100,
		Size:       0,
	}

	formatted := FormatResponse(resp)
	if !strings.Contains(formatted, "204") {
		t.Error("expected formatted response to contain status code")
	}
}

func TestUnitPrintResponse(t *testing.T) {
	resp := &types.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Headers:    make(map[string][]string),
		Body:       []byte("OK"),
		Duration:   time.Millisecond * 100,
		Size:       2,
	}

	PrintResponse(resp)
}

func TestUnitIsJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"valid JSON", []byte(`{"key": "value"}`), true},
		{"valid JSON array", []byte(`[1, 2, 3]`), true},
		{"invalid JSON", []byte("not json"), false},
		{"empty", []byte(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSON(tt.data)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

