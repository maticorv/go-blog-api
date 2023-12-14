package comments

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

// ICommentService define un servicio para obtener comments
type ICommentService interface {
	GetComments(title string, userID int, id int) (*[]data.Comment, error)
	GetComment(id int) (*data.Comment, error)
	CreateComment(album data.Comment) (*data.Comment, error)
	UpdateComment(id int, album data.Comment) (*data.Comment, error)
	DeleteComment(id int) error
}

// CommentService implementa el servicio utilizando JSONPlaceholder
type CommentService struct {
	restClient restclient.IRestClient
}

// GetComments obtiene comments desde JSONPlaceholder
func (s *CommentService) GetComments(title string, userID int, id int) (*[]data.Comment, error) {
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
	var comments []data.Comment
	err = json.NewDecoder(resp.Body).Decode(&comments)
	if err != nil {
		return nil, err
	}
	return &comments, nil
}

func (s *CommentService) GetComment(id int) (*data.Comment, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "GET"
	resp, err := s.restClient.NewRequest(method, url, bytes.NewBuffer(nil), nil)
	if err != nil {
		return nil, err
	}
	var album data.Comment
	err = json.NewDecoder(resp.Body).Decode(&album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (s *CommentService) CreateComment(album data.Comment) (*data.Comment, error) {
	validate := validator.New()
	if err := validate.Struct(album); err != nil {
		return nil, fmt.Errorf("Comment can´t be created. Status Code: %s", err)
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
		return nil, fmt.Errorf("Comment can´t be created. Status Code: %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Comment can´t be created. Status Code: %d", resp.StatusCode)
	}
	var createdComment data.Comment
	err = json.NewDecoder(resp.Body).Decode(&createdComment)
	if err != nil {
		return nil, err
	}
	return &createdComment, nil
}

func (s *CommentService) UpdateComment(id int, album data.Comment) (*data.Comment, error) {
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
		return nil, fmt.Errorf("Comment can´t be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedComment data.Comment
	err = json.NewDecoder(resp.Body).Decode(&updatedComment)
	if err != nil {
		return nil, err
	}
	return &updatedComment, nil
}

func (s *CommentService) PatchComment(id int, album data.Comment) (*data.Comment, error) {
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
		return nil, fmt.Errorf("Comment can't be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedComment data.Comment
	err = json.NewDecoder(resp.Body).Decode(&updatedComment)
	if err != nil {
		return nil, err
	}
	return &updatedComment, nil
}

func (s *CommentService) DeleteComment(id int) error {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "DELETE"
	payload := bytes.NewBuffer(nil)
	resp, err := s.restClient.NewRequest(method, url, payload, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Comment cannnt be created. Status Code: %d", resp.StatusCode)
	}
	return nil
}

// NewCommentCommentService crea una nueva instancia del servicio de comments
func NewCommentService(client restclient.IRestClient) ICommentService {
	return &CommentService{
		restClient: client,
	}
}
