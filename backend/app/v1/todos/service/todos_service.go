package todos

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

// ITodoService define un servicio para obtener todos
type ITodoService interface {
	GetTodos(title string, userID int, id int) (*[]data.Todo, error)
	GetTodo(id int) (*data.Todo, error)
	CreateTodo(album data.Todo) (*data.Todo, error)
	UpdateTodo(id int, album data.Todo) (*data.Todo, error)
	DeleteTodo(id int) error
}

// TodoService implementa el servicio utilizando JSONPlaceholder
type TodoService struct {
	restClient restclient.IRestClient
}

// GetTodos obtiene todos desde JSONPlaceholder
func (s *TodoService) GetTodos(title string, userID int, id int) (*[]data.Todo, error) {
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
	var todos []data.Todo
	err = json.NewDecoder(resp.Body).Decode(&todos)
	if err != nil {
		return nil, err
	}
	return &todos, nil
}

func (s *TodoService) GetTodo(id int) (*data.Todo, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "GET"
	resp, err := s.restClient.NewRequest(method, url, bytes.NewBuffer(nil), nil)
	if err != nil {
		return nil, err
	}
	var album data.Todo
	err = json.NewDecoder(resp.Body).Decode(&album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (s *TodoService) CreateTodo(album data.Todo) (*data.Todo, error) {
	validate := validator.New()
	if err := validate.Struct(album); err != nil {
		return nil, fmt.Errorf("Todo can´t be created. Status Code: %s", err)
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
		return nil, fmt.Errorf("Todo can´t be created. Status Code: %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Todo can´t be created. Status Code: %d", resp.StatusCode)
	}
	var createdTodo data.Todo
	err = json.NewDecoder(resp.Body).Decode(&createdTodo)
	if err != nil {
		return nil, err
	}
	return &createdTodo, nil
}

func (s *TodoService) UpdateTodo(id int, album data.Todo) (*data.Todo, error) {
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
		return nil, fmt.Errorf("Todo can´t be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedTodo data.Todo
	err = json.NewDecoder(resp.Body).Decode(&updatedTodo)
	if err != nil {
		return nil, err
	}
	return &updatedTodo, nil
}

func (s *TodoService) PatchTodo(id int, album data.Todo) (*data.Todo, error) {
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
		return nil, fmt.Errorf("Todo can't be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedTodo data.Todo
	err = json.NewDecoder(resp.Body).Decode(&updatedTodo)
	if err != nil {
		return nil, err
	}
	return &updatedTodo, nil
}

func (s *TodoService) DeleteTodo(id int) error {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "DELETE"
	payload := bytes.NewBuffer(nil)
	resp, err := s.restClient.NewRequest(method, url, payload, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Todo cannnt be created. Status Code: %d", resp.StatusCode)
	}
	return nil
}

// NewTodoTodoService crea una nueva instancia del servicio de todos
func NewTodoService(client restclient.IRestClient) ITodoService {
	return &TodoService{
		restClient: client,
	}
}
