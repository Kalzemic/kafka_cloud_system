package clients

import (
	"api/models"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type ConsumerClient interface {
	Poll(poll *models.PollRequest) ([]models.Post, error)
	Listen(ctx context.Context) (<-chan models.Post, error)
}

type APIConsumerClient struct {
	Client  *http.Client
	BaseURL string
}

func (cc *APIConsumerClient) Poll(poll *models.PollRequest) ([]models.Post, error) {
	data, err := json.Marshal(poll)
	if err != nil {
		return nil, ErrRequestFailure
	}
	resp, err := cc.Client.Post(cc.BaseURL+"/poll", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, ErrConsumerServiceUnavailable
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return nil, ErrInvalidInput
	case http.StatusOK:
		var posts []models.Post
		if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
			return nil, ErrConsumerServiceUnavailable
		}
		return posts, nil
	default:
		return nil, ErrConsumerServiceUnavailable
	}

}

func (cc *APIConsumerClient) Listen(ctx context.Context) (<-chan models.Post, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		cc.BaseURL+"/listen",
		nil,
	)
	if err != nil {
		return nil, ErrRequestFailure
	}

	req.Header.Set("Accept", "text/event-stream")

	resp, err := cc.Client.Do(req)
	if err != nil {
		return nil, ErrConsumerServiceUnavailable
	}

	out := make(chan models.Post)

	go func() {
		defer resp.Body.Close()
		defer close(out)
		scanner := bufio.NewScanner(resp.Body)
		var event string
		var data string

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}
			line := scanner.Text()
			if line == "" {
				if event == "post" && data != "" {
					var post models.Post
					if err := json.Unmarshal([]byte(data), &post); err == nil {
						out <- post
					}
				}
				event = ""
				data = ""
				continue
			}
			if strings.HasPrefix(line, "event:") {
				event = strings.TrimSpace(line[len("event:"):])
			}

			if strings.HasPrefix(line, "data:") {
				data = strings.TrimSpace(line[len("data:"):])
			}
		}
	}()
	return out, nil
}
