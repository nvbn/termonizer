package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/nvbn/termonizer/internal/model"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestGenerate_Success(t *testing.T) {
	mockResponses := []ollamaResponse{
		{Response: "This is line 1", Done: false},
		{Response: "This is line 2", Done: true},
	}
	var buf bytes.Buffer
	for _, resp := range mockResponses {
		line, _ := json.Marshal(resp)
		buf.Write(line)
		buf.WriteRune('\n')
	}

	client := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(&buf),
			}, nil
		},
	}

	ollama := NewClient(client, "http://example.com", "test-model")
	ctx := context.Background()
	ch := ollama.Generate(ctx, "test prompt")

	var results []model.AiResponse
	for res := range ch {
		results = append(results, res)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 responses, got %d", len(results))
	}
	if results[0].Text != "This is line 1" || results[1].Text != "This is line 2" {
		t.Errorf("Unexpected responses: %+v", results)
	}
}

func TestGenerate_InvalidJSON(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte("invalid json\n"))),
			}, nil
		},
	}

	ollama := NewClient(client, "http://example.com", "test-model")
	ctx := context.Background()
	ch := ollama.Generate(ctx, "test prompt")

	res, ok := <-ch
	if !ok {
		t.Fatal("Channel closed unexpectedly")
	}

	if res.Error == nil {
		t.Fatal("Expected an error due to invalid JSON, got none")
	}
}

func TestGenerate_HTTPError(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader(nil)),
			}, nil
		},
	}

	ollama := NewClient(client, "http://example.com", "test-model")
	ctx := context.Background()
	ch := ollama.Generate(ctx, "test prompt")

	res, ok := <-ch
	if !ok {
		t.Fatal("Channel closed unexpectedly")
	}

	if res.Error == nil {
		t.Fatal("Expected an error due to HTTP error, got none")
	}
}

func TestGenerate_EmptyStream(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(nil)),
			}, nil
		},
	}

	ollama := NewClient(client, "http://example.com", "test-model")
	ctx := context.Background()
	ch := ollama.Generate(ctx, "test prompt")

	// The channel should close without errors or responses
	res, ok := <-ch
	if ok {
		t.Errorf("Expected channel to close, but got response: %+v", res)
	}
}

func TestGenerate_ContextCanceled(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			time.Sleep(2 * time.Second)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(nil)),
			}, nil
		},
	}

	ollama := NewClient(client, "http://example.com", "test-model")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := ollama.Generate(ctx, "test prompt")

	res, ok := <-ch
	if ok {
		t.Errorf("Expected channel to close due to context cancelation, but got response: %+v", res)
	}
}
