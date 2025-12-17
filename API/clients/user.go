package clients

import (
	"api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type UserClient interface {
	CreateUser(user *models.UserRequest) (*models.UserResponse, error)
	UpdateUser(user *models.UserRequest) error
	GetbyEmailDomain(domain string, page int, size int) ([]models.UserResponse, error)
	FindUser(email string, password string) (*models.UserResponse, error)
	GetAllUsers(page int, size int) ([]models.UserResponse, error)
	GetUsersbyRoles(role string, page int, size int) ([]models.UserResponse, error)
	GetUsersbyRegistrationToday(page int, size int) ([]models.UserResponse, error)
	DeleteUsers() error
}

type APIUserClient struct {
	baseURL string
	client  *http.Client
}

func (uc *APIUserClient) CreateUser(user *models.UserRequest) (*models.UserResponse, error) {

	data, err := json.Marshal(user)
	if err != nil {
		return nil, ErrRequestFailure
	}
	resp, err := uc.client.Post(uc.baseURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, ErrUserServiceUnavailable
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:

		var created models.UserResponse
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			return nil, ErrUserServiceUnavailable
		}

		return &created, nil
	case http.StatusConflict:
		return nil, ErrUserAlreadyExists
	case http.StatusBadRequest:
		return nil, ErrInvalidInput
	default:
		return nil, ErrUserServiceUnavailable
	}

}

func (uc *APIUserClient) UpdateUser(user *models.UserRequest) error {

	data, err := json.Marshal(user)
	if err != nil {
		return ErrRequestFailure
	}

	req, err := http.NewRequest("PUT", uc.baseURL+"/"+url.PathEscape(user.Email)+"?password="+url.QueryEscape(user.Password), bytes.NewReader(data))
	if err != nil {
		return ErrRequestFailure
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		return ErrUserServiceUnavailable
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrUserNotFound
	case http.StatusOK:
		return nil
	default:
		return ErrUserServiceUnavailable
	}

}

func (uc *APIUserClient) FindUser(email string, password string) (*models.UserResponse, error) {
	resp, err := uc.client.Get(uc.baseURL + fmt.Sprintf("/%s?password=%s", url.PathEscape(email), url.QueryEscape(password)))
	if err != nil {
		return nil, ErrUserServiceUnavailable
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, ErrUserNotFound
	case http.StatusOK:
		var created models.UserResponse
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			return nil, ErrUserServiceUnavailable
		}

		return &created, nil
	default:
		return nil, ErrUserServiceUnavailable
	}
}

func (uc *APIUserClient) GetAllUsers(page int, size int) ([]models.UserResponse, error) {

	resp, err := uc.client.Get(uc.baseURL + fmt.Sprintf("?&page=%d&size=%d", page, size))
	if err != nil {
		return nil, ErrUserServiceUnavailable
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var users []models.UserResponse
		if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
			return nil, ErrUserServiceUnavailable
		}
		return users, nil
	default:
		return nil, ErrUserServiceUnavailable
	}

}

func (uc *APIUserClient) GetUsersbyRoles(role string, page int, size int) ([]models.UserResponse, error) {

	resp, err := uc.client.Get(uc.baseURL + fmt.Sprintf("?criteria=%s&value=%s&page=%d&size=%d", url.QueryEscape("byRole"), url.QueryEscape(role), page, size))
	if err != nil {
		return nil, ErrUserServiceUnavailable
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var users []models.UserResponse
		if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
			return nil, ErrUserServiceUnavailable
		}
		return users, nil
	default:
		return nil, ErrUserServiceUnavailable
	}

}

func (uc *APIUserClient) GetbyEmailDomain(domain string, page int, size int) ([]models.UserResponse, error) {

	resp, err := uc.client.Get(uc.baseURL + fmt.Sprintf("?criteria=%s&value=%s&page=%d&size=%d", url.QueryEscape("byEmailDomain"), url.QueryEscape(domain), page, size))
	if err != nil {
		return nil, ErrUserServiceUnavailable
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var users []models.UserResponse
		if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
			return nil, ErrUserServiceUnavailable
		}
		return users, nil
	default:
		return nil, ErrUserServiceUnavailable
	}

}

func (uc *APIUserClient) GetUsersbyRegistrationToday(page int, size int) ([]models.UserResponse, error) {

	resp, err := uc.client.Get(uc.baseURL + fmt.Sprintf("?criteria=%s&page=%d&size=%d", url.QueryEscape("byRegistrationToday"), page, size))
	if err != nil {
		return nil, ErrUserServiceUnavailable
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var users []models.UserResponse
		if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
			return nil, ErrUserServiceUnavailable
		}
		return users, nil
	default:
		return nil, ErrUserServiceUnavailable
	}

}

func (uc *APIUserClient) DeleteUsers() error {
	req, err := http.NewRequest("DELETE", uc.baseURL, nil)
	if err != nil {
		return ErrRequestFailure
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		return ErrSendingFailure
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return ErrUserServiceUnavailable
	}

}
