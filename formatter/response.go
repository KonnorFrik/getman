package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KonnorFrik/getman/types"
	"github.com/fatih/color"
)

func FormatResponse(resp *types.Response) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Status: %d %s\n", resp.StatusCode, resp.Status))
	sb.WriteString(fmt.Sprintf("Duration: %v\n", resp.Duration))
	sb.WriteString(fmt.Sprintf("Size: %d bytes\n", resp.Size))

	if len(resp.Headers) > 0 {
		sb.WriteString("\nHeaders:\n")
		for k, v := range resp.Headers {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ", ")))
		}
	}

	if len(resp.Body) > 0 {
		sb.WriteString("\nBody:\n")
		var bodyStr string
		if isJSON(resp.Body) {
			var jsonObj any
			if err := json.Unmarshal(resp.Body, &jsonObj); err == nil {
				prettyJSON, err := json.MarshalIndent(jsonObj, "", "  ")
				if err == nil {
					bodyStr = string(prettyJSON)
				} else {
					bodyStr = string(resp.Body)
				}
			} else {
				bodyStr = string(resp.Body)
			}
		} else {
			bodyStr = string(resp.Body)
		}
		sb.WriteString(bodyStr)
		sb.WriteString("\n")
	}

	return sb.String()
}

func PrintResponse(resp *types.Response) {
	var statusColor *color.Color
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		statusColor = color.New(color.FgGreen)
	case resp.StatusCode >= 300 && resp.StatusCode < 400:
		statusColor = color.New(color.FgYellow)
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		statusColor = color.New(color.FgRed)
	default:
		statusColor = color.New(color.FgMagenta)
	}

	statusColor.Printf("Status: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("Duration: %v\n", resp.Duration)
	fmt.Printf("Size: %d bytes\n", resp.Size)

	if len(resp.Headers) > 0 {
		fmt.Println("\nHeaders:")
		for k, v := range resp.Headers {
			fmt.Printf("  %s: %s\n", k, strings.Join(v, ", "))
		}
	}

	if len(resp.Body) > 0 {
		fmt.Println("\nBody:")
		var bodyStr string
		if isJSON(resp.Body) {
			var jsonObj any
			if err := json.Unmarshal(resp.Body, &jsonObj); err == nil {
				prettyJSON, err := json.MarshalIndent(jsonObj, "", "  ")
				if err == nil {
					bodyStr = string(prettyJSON)
				} else {
					bodyStr = string(resp.Body)
				}
			} else {
				bodyStr = string(resp.Body)
			}
		} else {
			bodyStr = string(resp.Body)
		}
		fmt.Println(bodyStr)
	}
}

func isJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}
