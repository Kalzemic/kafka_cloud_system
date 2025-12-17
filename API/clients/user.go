package clients

import (
	"api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserClient interface {
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetbyEmailDomain(domain string, page int, size int) ([]models.User, error)
	FindUser(email string, password string) (*models.User, error)
	GetAllUsers(page int, size int) ([]models.User, error)
	GetUsersbyRoles(role string, page int, limit int) ([]models.User, error)
	GetUsersbyRegistrationToday(page int, size int) ([]models.User, error)
	DeleteUsers() error
}

type APIUserClient struct {
	baseURL string
	client  *http.Client
}

func (uc *APIUserClient) CreateUser(user *models.User) error {

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed converting user to json %w", err)
	}
	resp, err := uc.client.Post(uc.baseURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("from user service %d : %s", resp.StatusCode, string(body))
	}

	return nil
}
