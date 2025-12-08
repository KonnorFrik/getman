/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package http_server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	testServer    *http.Server
	testServerURL string
	testServerMu  sync.Mutex
	listener      net.Listener
)

type TestServer struct {
	server *http.Server
	url    string
}

func StartTestServer() (*TestServer, error) {
	testServerMu.Lock()
	defer testServerMu.Unlock()

	if testServer != nil {
		return &TestServer{
			server: testServer,
			url:    testServerURL,
		}, nil
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/echo", handleEcho)
	mux.HandleFunc("/status/", handleStatus)
	mux.HandleFunc("/headers", handleHeaders)
	mux.HandleFunc("/auth/basic", handleBasicAuth)
	mux.HandleFunc("/auth/bearer", handleBearerAuth)
	mux.HandleFunc("/auth/apikey", handleAPIKeyAuth)
	mux.HandleFunc("/cookies", handleCookies)
	mux.HandleFunc("/body", handleBody)
	mux.HandleFunc("/delay/", handleDelay)

	var err error
	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	server := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("failed to start test server: %v", err))
		}
	}()

	time.Sleep(100 * time.Millisecond)

	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", port)
	testServer = server
	testServerURL = url

	return &TestServer{
		server: server,
		url:    url,
	}, nil
}

func StopTestServer() error {
	testServerMu.Lock()
	defer testServerMu.Unlock()

	if testServer == nil {
		return nil
	}

	err := testServer.Close()
	if listener != nil {
		listener.Close()
		listener = nil
	}
	testServer = nil
	testServerURL = ""
	return err
}

func GetServerURL() string {
	testServerMu.Lock()
	defer testServerMu.Unlock()
	return testServerURL
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"method":  r.Method,
		"path":    r.URL.Path,
		"headers": r.Header,
		"query":   r.URL.Query(),
	}

	if r.Body != nil {
		body := make([]byte, 0)
		buf := make([]byte, 1024)
		for {
			n, err := r.Body.Read(buf)
			if n > 0 {
				body = append(body, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
		if len(body) > 0 {
			response["body"] = string(body)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("Status: %d", code)))
}

func handleHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Header)
}

func handleBasicAuth(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != "testuser" || password != "testpass" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authenticated"))
}

func handleBearerAuth(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != "testtoken" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authenticated"))
}

func handleAPIKeyAuth(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = r.URL.Query().Get("api_key")
	}

	if apiKey != "testapikey" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authenticated"))
}

func handleCookies(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "testcookie",
		Value:   "testvalue",
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	response := map[string]interface{}{
		"cookies": r.Cookies(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleBody(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := r.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	response := map[string]interface{}{
		"body":    string(body),
		"headers": r.Header,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleDelay(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	time.Sleep(time.Duration(seconds) * time.Second)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delayed response"))
}
