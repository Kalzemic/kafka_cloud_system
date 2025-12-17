package clients

import (
	"api/models"
	"bytes"
	"encoding/json"
	"net/http"
)

type ConsumerClient interface {
	Poll(poll *models.PollRequest) ([]models.Post, error)
}

type APIConsumerClient struct {
	client  *http.Client
	baseURL string
}

func (cc *APIConsumerClient) Poll(poll *models.PollRequest) ([]models.Post, error) {
	data, err := json.Marshal(poll)
	if err != nil {
		return nil, ErrRequestFailure
	}
	resp, err := cc.client.Post(cc.baseURL+"/poll", "application/json", bytes.NewReader(data))
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
