/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package types

import "time"

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    *RequestBody      `json:"body,omitempty"`
	Auth    *Auth             `json:"auth,omitempty"`
	Timeout *Timeout          `json:"timeout,omitempty"`
	Cookies *CookieSettings   `json:"cookies,omitempty"`
}

type RequestBody struct {
	Type        string `json:"type"`
	Content     []byte `json:"content"`
	ContentType string `json:"content_type,omitempty"`
}

type Auth struct {
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
	APIKey   string `json:"api_key,omitempty"`
	KeyName  string `json:"key_name,omitempty"`
	Location string `json:"location,omitempty"`
}

type Timeout struct {
	Connect time.Duration `json:"connect"`
	Read    time.Duration `json:"read"`
}

type CookieSettings struct {
	AutoManage bool `json:"auto_manage"`
}

type Response struct {
	StatusCode int                 `json:"status_code"`
	Status     string              `json:"status"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"body"`
	Duration   time.Duration       `json:"duration"`
	Size       int64               `json:"size"`
}

type RequestItem struct {
	Name    string   `json:"name"`
	Request *Request `json:"request"`
}

type RequestExecution struct {
	Request   *Request      `json:"request"`
	Response  *Response     `json:"response,omitempty"`
	Error     string        `json:"error,omitempty"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

type ExecutionResult struct {
	CollectionName string              `json:"collection_name"`
	Environment    string              `json:"environment"`
	StartTime      time.Time           `json:"start_time"`
	EndTime        time.Time           `json:"end_time"`
	TotalDuration  time.Duration       `json:"total_duration"`
	Requests       []*RequestExecution `json:"requests"`
	Statistics     *Statistics         `json:"statistics"`
}

type Statistics struct {
	Total   int           `json:"total"`
	Success int           `json:"success"`
	Failed  int           `json:"failed"`
	AvgTime time.Duration `json:"avg_time"`
	MinTime time.Duration `json:"min_time"`
	MaxTime time.Duration `json:"max_time"`
}

type LogEntry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}
