package testkit

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// RegistryResponse controls how a RegistryStub answers every request. A zero
// status means 200 OK. Hang waits until the caller cancels or the stub closes,
// which makes timeout tests deterministic without a long sleep.
type RegistryResponse struct {
	StatusCode  int
	ContentType string
	Body        []byte
	Delay       time.Duration
	Hang        bool
}

// RegistryStub is a local HTTP server whose response can be changed during a
// test. It is safe to use from concurrent request handlers.
type RegistryStub struct {
	server *httptest.Server

	mu       sync.RWMutex
	response RegistryResponse
	requests []*http.Request
	done     chan struct{}
	close    sync.Once
}

// NewRegistryStub starts a controllable local registry. Callers normally need
// only one line, for example:
//
//	registry := testkit.NewRegistryStub(t, testkit.RegistryResponse{StatusCode: 500})
func NewRegistryStub(t testing.TB, response RegistryResponse) *RegistryStub {
	t.Helper()

	stub := &RegistryStub{response: cloneResponse(response), done: make(chan struct{})}
	stub.server = httptest.NewServer(http.HandlerFunc(stub.serveHTTP))
	t.Cleanup(stub.Close)
	return stub
}

// URL returns the base URL to inject into application configuration.
func (s *RegistryStub) URL() string {
	return s.server.URL
}

// Client returns an HTTP client configured for this test server.
func (s *RegistryStub) Client() *http.Client {
	return s.server.Client()
}

// SetResponse changes the behaviour used by subsequent requests.
func (s *RegistryStub) SetResponse(response RegistryResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.response = cloneResponse(response)
}

// RequestCount returns the number of requests observed so far.
func (s *RegistryStub) RequestCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.requests)
}

// Close stops the server. It is safe to call more than once.
func (s *RegistryStub) Close() {
	s.close.Do(func() {
		close(s.done)
		s.server.Close()
	})
}

func (s *RegistryStub) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.requests = append(s.requests, r.Clone(r.Context()))
	response := cloneResponse(s.response)
	s.mu.Unlock()

	if response.Hang {
		select {
		case <-r.Context().Done():
		case <-s.done:
		}
		return
	}
	if response.Delay > 0 {
		timer := time.NewTimer(response.Delay)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-r.Context().Done():
			return
		case <-s.done:
			return
		}
	}

	status := response.StatusCode
	if status == 0 {
		status = http.StatusOK
	}
	if response.ContentType != "" {
		w.Header().Set("Content-Type", response.ContentType)
	}
	w.WriteHeader(status)
	_, _ = w.Write(response.Body)
}

func cloneResponse(response RegistryResponse) RegistryResponse {
	response.Body = append([]byte(nil), response.Body...)
	return response
}
