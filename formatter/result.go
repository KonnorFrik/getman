package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/KonnorFrik/getman/types"
	"github.com/fatih/color"
)

var (
	colorFgMagneta = color.New(color.FgHiMagenta)
	colorFgCyan = color.New(color.FgHiCyan)
)

func FormatRequest(req *types.Request) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s\n", req.Method, req.URL))

	if len(req.Headers) > 0 {
		sb.WriteString("\nHeaders:\n")
		for k, v := range req.Headers {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	if req.Auth != nil {
		sb.WriteString("\nAuth:\n")
		sb.WriteString(fmt.Sprintf("  Type: %s\n", req.Auth.Type))
		if req.Auth.Username != "" {
			sb.WriteString(fmt.Sprintf("  Username: %s\n", req.Auth.Username))
		}
		if req.Auth.Token != "" {
			sb.WriteString(fmt.Sprintf("  Token: %s\n", maskToken(req.Auth.Token)))
		}
		if req.Auth.APIKey != "" {
			sb.WriteString(fmt.Sprintf("  APIKey: %s\n", maskToken(req.Auth.APIKey)))
		}
	}

	if req.Body != nil && len(req.Body.Content) > 0 {
		sb.WriteString("\nBody:\n")
		var bodyStr string
		if isJSON(req.Body.Content) {
			var jsonObj interface{}
			if err := json.Unmarshal(req.Body.Content, &jsonObj); err == nil {
				prettyJSON, err := json.MarshalIndent(jsonObj, "", "  ")
				if err == nil {
					bodyStr = string(prettyJSON)
				} else {
					bodyStr = string(req.Body.Content)
				}
			} else {
				bodyStr = string(req.Body.Content)
			}
		} else {
			bodyStr = string(req.Body.Content)
		}
		sb.WriteString(bodyStr)
		sb.WriteString("\n")
	}

	return sb.String()
}

func PrintRequest(req *types.Request) {
	fmt.Printf("%s %s\n", req.Method, req.URL)

	if len(req.Headers) > 0 {
		fmt.Println("\nHeaders:")
		for k, v := range req.Headers {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}

	if req.Auth != nil {
		fmt.Println("\nAuth:")
		fmt.Printf("  Type: %s\n", req.Auth.Type)
		if req.Auth.Username != "" {
			fmt.Printf("  Username: %s\n", req.Auth.Username)
		}
		if req.Auth.Token != "" {
			fmt.Printf("  Token: %s\n", maskToken(req.Auth.Token))
		}
		if req.Auth.APIKey != "" {
			fmt.Printf("  APIKey: %s\n", maskToken(req.Auth.APIKey))
		}
	}

	if req.Body != nil && len(req.Body.Content) > 0 {
		fmt.Println("\nBody:")
		var bodyStr string
		if isJSON(req.Body.Content) {
			var jsonObj interface{}
			if err := json.Unmarshal(req.Body.Content, &jsonObj); err == nil {
				prettyJSON, err := json.MarshalIndent(jsonObj, "", "  ")
				if err == nil {
					bodyStr = string(prettyJSON)
				} else {
					bodyStr = string(req.Body.Content)
				}
			} else {
				bodyStr = string(req.Body.Content)
			}
		} else {
			bodyStr = string(req.Body.Content)
		}
		fmt.Println(bodyStr)
	}
}

func FormatExecutionResult(result *types.ExecutionResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Collection: %s\n", result.CollectionName))
	sb.WriteString(fmt.Sprintf("Environment: %s\n", result.Environment))
	sb.WriteString(fmt.Sprintf("Start Time: %s\n", result.StartTime.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("End Time: %s\n", result.EndTime.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Total Duration: %v\n", result.TotalDuration))

	if result.Statistics != nil {
		sb.WriteString("\nStatistics:\n")
		sb.WriteString(FormatStatistics(result.Statistics))
	}

	sb.WriteString("\nRequests:\n")
	for i, req := range result.Requests {
		sb.WriteString(fmt.Sprintf("\n%d. %s %s\n", i+1, req.Request.Method, req.Request.URL))
		if req.Error != "" {
			sb.WriteString(fmt.Sprintf("   Error: %s\n", req.Error))
		} else if req.Response != nil {
			sb.WriteString(fmt.Sprintf("   Status: %d\n", req.Response.StatusCode))
			sb.WriteString(fmt.Sprintf("   Duration: %v\n", req.Duration))
		}
	}

	return sb.String()
}

func PrintExecutionResult(result *types.ExecutionResult) {
	fmt.Printf("Collection: %s\n", result.CollectionName)
	fmt.Printf("Environment: %s\n", result.Environment)
	fmt.Printf("Start Time: %s\n", result.StartTime.Format(time.RFC3339))
	fmt.Printf("End Time: %s\n", result.EndTime.Format(time.RFC3339))
	fmt.Printf("Total Duration: %v\n", result.TotalDuration)

	if result.Statistics != nil {
		fmt.Println("\nStatistics:")
		PrintStatistics(result.Statistics)
	}

	fmt.Println("\nRequests:")

	for i, req := range result.Requests {
		fmt.Printf("\n%d. %s %s\n", i+1, req.Request.Method, req.Request.URL)

		if req.Error != "" {
			color.Red("   Error: %s\n", req.Error)

		} else if req.Response != nil {
			var statusColor *color.Color

			switch {
			case req.Response.StatusCode >= 200 && req.Response.StatusCode < 300:
				statusColor = color.New(color.FgGreen)

			case req.Response.StatusCode >= 300 && req.Response.StatusCode < 400:
				statusColor = color.New(color.FgYellow)

			case req.Response.StatusCode >= 400 && req.Response.StatusCode < 500:
				statusColor = color.New(color.FgRed)

			default:
				statusColor = color.New(color.FgMagenta)
			}

			statusColor.Printf("   Status: %d\n", req.Response.StatusCode)
			fmt.Printf("   Duration: %v\n", req.Duration)
			colorFgMagneta.Printf("   Headers:\n")

			for k, v := range req.Response.Headers {
				colorFgCyan.Printf("\t%s", k)
				fmt.Printf(": %s\n", strings.Join(v, " "))
			}

			colorFgMagneta.Printf("   Body:\n")
			fmt.Printf("%s\n", string(req.Response.Body))
		}
	}
}

func FormatStatistics(stats *types.Statistics) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  Total: %d\n", stats.Total))
	sb.WriteString(fmt.Sprintf("  Success: %d\n", stats.Success))
	sb.WriteString(fmt.Sprintf("  Failed: %d\n", stats.Failed))
	sb.WriteString(fmt.Sprintf("  Avg Time: %v\n", stats.AvgTime))
	sb.WriteString(fmt.Sprintf("  Min Time: %v\n", stats.MinTime))
	sb.WriteString(fmt.Sprintf("  Max Time: %v\n", stats.MaxTime))
	return sb.String()
}

func PrintStatistics(stats *types.Statistics) {
	fmt.Printf("  Total: %d\n", stats.Total)
	color.Green("  Success: %d\n", stats.Success)
	color.Red("  Failed: %d\n", stats.Failed)
	fmt.Printf("  Avg Time: %v\n", stats.AvgTime)
	fmt.Printf("  Min Time: %v\n", stats.MinTime)
	fmt.Printf("  Max Time: %v\n", stats.MaxTime)
}

func maskToken(token string) string {
	if len(token) < 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
