package posts

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

// PostPostService define un servicio para obtener posts
type IPostService interface {
	GetPosts(title string, userID int, id int) (*[]data.Post, error)
	GetPost(id int) (*data.Post, error)
	CreatePost(post data.Post) (*data.Post, error)
	UpdatePost(id int, post data.Post) (*data.Post, error)
	DeletePost(id int) error
}

// PostService implementa el servicio utilizando JSONPlaceholder
type PostService struct {
	restClient restclient.IRestClient
}

// GetPosts obtiene posts desde JSONPlaceholder
func (s *PostService) GetPosts(title string, userID int, id int) (*[]data.Post, error) {
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
	var posts []data.Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (s *PostService) GetPost(id int) (*data.Post, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "GET"
	resp, err := s.restClient.NewRequest(method, url, bytes.NewBuffer(nil), nil)
	if err != nil {
		return nil, err
	}
	var post data.Post
	err = json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) CreatePost(post data.Post) (*data.Post, error) {
	validate := validator.New()
	if err := validate.Struct(post); err != nil {
		return nil, fmt.Errorf("Post can´t be created. Status Code: %s", err)
	}

	url := baseUrl
	method := "POST"
	body, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(body)
	headers := map[string]string{"Content-Type": " application/json"}
	resp, err := s.restClient.NewRequest(method, url, payload, headers)
	if err != nil {
		return nil, fmt.Errorf("Post can´t be created. Status Code: %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Post can´t be created. Status Code: %d", resp.StatusCode)
	}
	var createdPost data.Post
	err = json.NewDecoder(resp.Body).Decode(&createdPost)
	if err != nil {
		return nil, err
	}
	return &createdPost, nil
}

func (s *PostService) UpdatePost(id int, post data.Post) (*data.Post, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "PUT"
	body, err := json.Marshal(post)
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
		return nil, fmt.Errorf("Post can´t be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedPost data.Post
	err = json.NewDecoder(resp.Body).Decode(&updatedPost)
	if err != nil {
		return nil, err
	}
	return &updatedPost, nil
}

func (s *PostService) PatchPost(id int, post data.Post) (*data.Post, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "PATCH"
	body, err := json.Marshal(post)
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
		return nil, fmt.Errorf("Post can't be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedPost data.Post
	err = json.NewDecoder(resp.Body).Decode(&updatedPost)
	if err != nil {
		return nil, err
	}
	return &updatedPost, nil
}

func (s *PostService) DeletePost(id int) error {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "DELETE"
	payload := bytes.NewBuffer(nil)
	resp, err := s.restClient.NewRequest(method, url, payload, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Post cannnt be created. Status Code: %d", resp.StatusCode)
	}
	return nil
}

// NewPostPostService crea una nueva instancia del servicio de posts
func NewPostService(client restclient.IRestClient) IPostService {
	return &PostService{
		restClient: client,
	}
}
