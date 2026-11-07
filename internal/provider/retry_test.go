package provider

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompleteRetriesTransientErrors(t *testing.T) {
	for _, status := range []int{408, 429, 503} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if calls == 1 {
					w.WriteHeader(status)
					return
				}
				w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Done"}}]}`))
			}))
			defer server.Close()
			client := &Client{URL: server.URL, Key: "test", HTTP: server.Client()}
			message, _, err := client.Complete(context.Background(), Request{})
			if err != nil || message.Content != "Done" || calls != 2 {
				t.Fatalf("Unexpected retry result: %+v, %v, calls=%d", message, err, calls)
			}
		})
	}
}

func TestCompleteDoesNotRetryPermanentErrors(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	client := &Client{URL: server.URL, Key: "test", HTTP: server.Client()}
	if _, _, err := client.Complete(context.Background(), Request{}); err == nil || calls != 1 {
		t.Fatalf("Expected one failed attempt, got %d, %v", calls, err)
	}
}

func TestCompleteCancelsDuringRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusServiceUnavailable)
		cancel()
	}))
	defer server.Close()
	client := &Client{URL: server.URL, Key: "test", HTTP: server.Client()}
	if _, _, err := client.Complete(ctx, Request{}); !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatalf("Expected cancellation after one attempt, got %d, %v", calls, err)
	}
}

func TestDefaultClientHasNoTimeout(t *testing.T) {
	if client := FromEnv(); client.HTTP.Timeout != 0 {
		t.Fatalf("Unexpected provider timeout: %v", client.HTTP.Timeout)
	}
}
