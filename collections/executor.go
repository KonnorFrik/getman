/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package collections

import (
	"fmt"
	"time"

	"github.com/KonnorFrik/getman/core"
	"github.com/KonnorFrik/getman/types"
)

// CollectionExecutor executes collections of HTTP requests.
type CollectionExecutor struct {
	httpClient       *core.HTTPClient
	variableResolver *core.VariableResolver
}

// NewCollectionExecutor creates a new CollectionExecutor instance.
func NewCollectionExecutor(httpClient *core.HTTPClient, variableResolver *core.VariableResolver) *CollectionExecutor {
	return &CollectionExecutor{
		httpClient:       httpClient,
		variableResolver: variableResolver,
	}
}

// ExecuteCollection executes all requests in a collection.
func (ce *CollectionExecutor) ExecuteCollection(collection *Collection, environment string) (*types.ExecutionResult, error) {
	return ce.ExecuteCollectionSelective(collection, environment, nil)
}

// ExecuteCollectionAsync executes all requests in a collection asynchronously.
// It returns a channel that receives RequestExecution results as they complete.
// The channel is buffered with capacity 1. Results are sent in the order requests are processed.
// Note: The channel is not closed automatically after all requests complete.
func (ce *CollectionExecutor) ExecuteCollectionAsync(collection *Collection, environment string) <-chan *types.RequestExecution {
	ch := make(chan *types.RequestExecution, 1)

	go func() {
		defer close(ch)
		itemsToExecute := collection.Items

		for _, item := range itemsToExecute {
			req := item.Request
			resolvedReq, err := ce.resolveRequest(req)

			if err != nil {
				execution := &types.RequestExecution{
					Request:   req,
					Error:     fmt.Sprintf("failed to resolve variables: %v", err),
					Duration:  0,
					Timestamp: time.Now(),
				}
				ch <- execution
				continue
			}

			execStartTime := time.Now()
			response, err := ce.httpClient.Execute(resolvedReq)
			execDuration := time.Since(execStartTime)
			execution := &types.RequestExecution{
				Request:   resolvedReq,
				Duration:  execDuration,
				Timestamp: time.Now(),
			}

			if err != nil {
				execution.Error = err.Error()

			} else {
				execution.Response = response
			}

			ch <- execution
		}
	}()

	return ch
}

// ExecuteCollectionSelective executes only the specified requests from a collection.
func (ce *CollectionExecutor) ExecuteCollectionSelective(collection *Collection, environment string, itemNames []string) (*types.ExecutionResult, error) {
	startTime := time.Now()

	var itemsToExecute []*types.RequestItem
	if len(itemNames) == 0 {
		itemsToExecute = collection.Items

	} else {
		itemMap := make(map[string]*types.RequestItem)
		for _, item := range collection.Items {
			itemMap[item.Name] = item
		}

		for _, name := range itemNames {
			if item, ok := itemMap[name]; ok {
				itemsToExecute = append(itemsToExecute, item)
			}
		}
	}

	var (
		executions                []*types.RequestExecution
		totalDuration             time.Duration
		successCount, failedCount int
		minTime, maxTime          time.Duration
		firstTime                 = true
	)

	for _, item := range itemsToExecute {
		req := item.Request

		resolvedReq, err := ce.resolveRequest(req)
		if err != nil {
			execution := &types.RequestExecution{
				Request:   req,
				Error:     fmt.Sprintf("failed to resolve variables: %v", err),
				Duration:  0,
				Timestamp: time.Now(),
			}
			executions = append(executions, execution)
			failedCount++
			continue
		}

		execStartTime := time.Now()
		response, err := ce.httpClient.Execute(resolvedReq)
		execDuration := time.Since(execStartTime)
		totalDuration += execDuration

		if firstTime {
			minTime = execDuration
			maxTime = execDuration
			firstTime = false
		} else {
			if execDuration < minTime {
				minTime = execDuration
			}
			if execDuration > maxTime {
				maxTime = execDuration
			}
		}

		execution := &types.RequestExecution{
			Request:   resolvedReq,
			Duration:  execDuration,
			Timestamp: time.Now(),
		}

		if err != nil {
			execution.Error = err.Error()
			failedCount++
		} else {
			execution.Response = response
			if response.StatusCode >= 200 && response.StatusCode < 300 {
				successCount++
			} else {
				failedCount++
			}
		}

		executions = append(executions, execution)
	}

	endTime := time.Now()
	avgTime := time.Duration(0)
	if len(executions) > 0 {
		avgTime = totalDuration / time.Duration(len(executions))
	}

	result := &types.ExecutionResult{
		CollectionName: collection.Name,
		Environment:    environment,
		StartTime:      startTime,
		EndTime:        endTime,
		TotalDuration:  endTime.Sub(startTime),
		Requests:       executions,
		Statistics: &types.Statistics{
			Total:   len(executions),
			Success: successCount,
			Failed:  failedCount,
			AvgTime: avgTime,
			MinTime: minTime,
			MaxTime: maxTime,
		},
	}

	return result, nil
}

func (ce *CollectionExecutor) resolveRequest(req *types.Request) (*types.Request, error) {
	resolvedURL, err := ce.variableResolver.Resolve(req.URL)
	if err != nil {
		return nil, err
	}

	resolvedHeaders, err := ce.variableResolver.ResolveMap(req.Headers)
	if err != nil {
		return nil, err
	}

	resolvedReq := &types.Request{
		Method:  req.Method,
		URL:     resolvedURL,
		Headers: resolvedHeaders,
		Body:    req.Body,
		Auth:    req.Auth,
		Timeout: req.Timeout,
		Cookies: req.Cookies,
	}

	if req.Body != nil && len(req.Body.Content) > 0 {
		resolvedBodyContent, err := ce.variableResolver.Resolve(string(req.Body.Content))
		if err != nil {
			return nil, err
		}
		resolvedReq.Body = &types.RequestBody{
			Type:        req.Body.Type,
			Content:     []byte(resolvedBodyContent),
			ContentType: req.Body.ContentType,
		}
	}

	if req.Auth != nil {
		resolvedAuth := &types.Auth{
			Type:     req.Auth.Type,
			Username: req.Auth.Username,
			Password: req.Auth.Password,
			Token:    req.Auth.Token,
			APIKey:   req.Auth.APIKey,
			KeyName:  req.Auth.KeyName,
			Location: req.Auth.Location,
		}

		if req.Auth.Username != "" {
			resolvedAuth.Username, err = ce.variableResolver.Resolve(req.Auth.Username)
			if err != nil {
				return nil, err
			}
		}
		if req.Auth.Password != "" {
			resolvedAuth.Password, err = ce.variableResolver.Resolve(req.Auth.Password)
			if err != nil {
				return nil, err
			}
		}
		if req.Auth.Token != "" {
			resolvedAuth.Token, err = ce.variableResolver.Resolve(req.Auth.Token)
			if err != nil {
				return nil, err
			}
		}
		if req.Auth.APIKey != "" {
			resolvedAuth.APIKey, err = ce.variableResolver.Resolve(req.Auth.APIKey)
			if err != nil {
				return nil, err
			}
		}

		resolvedReq.Auth = resolvedAuth
	}

	return resolvedReq, nil
}
