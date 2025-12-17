package clients

import (
	"api/models"
	"bytes"
	"encoding/json"
	"net/http"
)

type ProducerClient interface {
	CreatePost(post *models.Post) error
}

type APIProducerClient struct {
	Client  *http.Client
	baseURL string
}

func (pc *APIProducerClient) CreatePost(post *models.Post) error {

	data, err := json.Marshal(post)
	if err != nil {
		return ErrRequestFailure
	}
	resp, err := pc.Client.Post(pc.baseURL+"/produce", "application/json", bytes.NewReader(data))
	if err != nil {
		return ErrProducerServiceUnavailable
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return ErrInvalidInput
	case http.StatusOK:
		return nil
	default:
		return ErrProducerServiceUnavailable
	}

}
