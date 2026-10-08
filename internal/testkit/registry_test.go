package testkit

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRegistryStubCanReturnConfiguredFailures(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusNotFound, http.StatusInternalServerError} {
		stub := NewRegistryStub(t, RegistryResponse{StatusCode: status})
		response, err := stub.Client().Get(stub.URL() + "/transactions/example.json")
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode != status {
			t.Fatalf("status = %d, want %d", response.StatusCode, status)
		}
	}

	stub := NewRegistryStub(t, RegistryResponse{})
	stub.SetResponse(RegistryResponse{
		StatusCode:  http.StatusOK,
		ContentType: "application/json",
		Body:        []byte(`{"not-valid-json"`),
	})
	response, err := stub.Client().Get(stub.URL() + "/transactions/broken.json")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `{"not-valid-json"` {
		t.Fatalf("body = %q", body)
	}
	if stub.RequestCount() != 1 {
		t.Fatalf("RequestCount = %d, want 1", stub.RequestCount())
	}
}

func TestRegistryStubCanHangUntilRequestIsCancelled(t *testing.T) {
	t.Parallel()

	stub := NewRegistryStub(t, RegistryResponse{Hang: true})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, stub.URL()+"/transactions/slow.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = stub.Client().Do(request)
	if err == nil {
		t.Fatal("request unexpectedly succeeded")
	}
	if ctx.Err() != context.DeadlineExceeded {
		t.Fatalf("context error = %v, want deadline exceeded", ctx.Err())
	}
}
