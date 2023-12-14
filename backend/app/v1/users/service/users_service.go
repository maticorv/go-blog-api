package users

import (
	"blog-api/app/clients/restclient"
	data "blog-api/data"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var baseUrl = os.Getenv("baseURL")

// IUserService define un servicio para obtener users
type IUserService interface {
	GetUsers(title string, userID int, id int) (*[]data.User, error)
	GetUser(id int) (*data.User, error)
	CreateUser(album data.User) (*data.User, error)
	UpdateUser(id int, album data.User) (*data.User, error)
	DeleteUser(id int) error
}

// UserService implementa el servicio utilizando JSONPlaceholder
type UserService struct {
	restClient restclient.IRestClient
}

// GetUsers obtiene users desde JSONPlaceholder
func (s *UserService) GetUsers(title string, userID int, id int) (*[]data.User, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if title != "" {
		q.Set("title", title)
	}
	if userID != 0 {
		q.Set("userId", strconv.Itoa(userID))
	}
	if id != 0 {
		q.Set("id", strconv.Itoa(id))
	}
	u.RawQuery = q.Encode()
	url := u.String()
	payload := bytes.NewBuffer(nil)
	method := "GET"
	resp, err := s.restClient.NewRequest(method, url, payload, nil)
	var users []data.User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (s *UserService) GetUser(id int) (*data.User, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "GET"
	resp, err := s.restClient.NewRequest(method, url, bytes.NewBuffer(nil), nil)
	if err != nil {
		return nil, err
	}
	var album data.User
	err = json.NewDecoder(resp.Body).Decode(&album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (s *UserService) CreateUser(album data.User) (*data.User, error) {
	validate := validator.New()
	if err := validate.Struct(album); err != nil {
		return nil, fmt.Errorf("User can´t be created. Status Code: %s", err)
	}

	url := baseUrl
	method := "POST"
	body, err := json.Marshal(album)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(body)
	headers := map[string]string{"Content-Type": " application/json"}
	resp, err := s.restClient.NewRequest(method, url, payload, headers)
	if err != nil {
		return nil, fmt.Errorf("User can´t be created. Status Code: %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("User can´t be created. Status Code: %d", resp.StatusCode)
	}
	var createdUser data.User
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

func (s *UserService) UpdateUser(id int, album data.User) (*data.User, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "PUT"
	body, err := json.Marshal(album)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(body)
	headers := map[string]string{"Content-Type": " application/json"}
	resp, err := s.restClient.NewRequest(method, url, payload, headers)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User can´t be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedUser data.User
	err = json.NewDecoder(resp.Body).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (s *UserService) PatchUser(id int, album data.User) (*data.User, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "PATCH"
	body, err := json.Marshal(album)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{"Content-Type": " application/json"}
	payload := bytes.NewBuffer(body)
	resp, err := s.restClient.NewRequest(method, url, payload, headers)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User can't be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedUser data.User
	err = json.NewDecoder(resp.Body).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (s *UserService) DeleteUser(id int) error {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "DELETE"
	payload := bytes.NewBuffer(nil)
	resp, err := s.restClient.NewRequest(method, url, payload, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("User cannnt be created. Status Code: %d", resp.StatusCode)
	}
	return nil
}

// NewUserUserService crea una nueva instancia del servicio de users
func NewUserService(client restclient.IRestClient) IUserService {
	return &UserService{
		restClient: client,
	}
}
