package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nvbn/termonizer/internal/model"
	"log"
	"net/http"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Ollama struct {
	client httpClient
	url    string
	model  string
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

func NewClient(client httpClient, url string, model string) *Ollama {
	return &Ollama{
		client: client,
		url:    url,
		model:  model,
	}
}

func (o *Ollama) Generate(ctx context.Context, prompt string) chan model.AiResponse {
	ch := make(chan model.AiResponse)

	go func() {
		defer close(ch)

		log.Printf("Generating with ollama: %s %v", prompt, o)

		reqData := ollamaRequest{
			Model:  o.model,
			Prompt: prompt,
			Stream: true,
		}

		reqDataJson, err := json.Marshal(reqData)
		if err != nil {
			ch <- model.AiResponse{Error: fmt.Errorf("unable to serialize ollama reqeust: %w", err)}
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", o.url, bytes.NewBuffer(reqDataJson))
		if err != nil {
			ch <- model.AiResponse{Error: fmt.Errorf("unable to create ollama request: %w", err)}
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := o.client.Do(req)
		if resp.StatusCode != http.StatusOK {
			ch <- model.AiResponse{Error: fmt.Errorf("ollama request failed: %w", err)}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			var respData ollamaResponse
			if err := json.Unmarshal([]byte(line), &respData); err != nil {
				ch <- model.AiResponse{Error: fmt.Errorf("error decoding ollama response line: %w", err)}
				return
			}

			select {
			case ch <- model.AiResponse{Text: respData.Response}:
				break
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch
}
