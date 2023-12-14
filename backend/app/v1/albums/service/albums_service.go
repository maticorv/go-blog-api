package albums

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

// AlbumAlbumService define un servicio para obtener albums
type IAlbumService interface {
	GetAlbums(title string, userID int, id int) (*[]data.Album, error)
	GetAlbum(id int) (*data.Album, error)
	CreateAlbum(album data.Album) (*data.Album, error)
	UpdateAlbum(id int, album data.Album) (*data.Album, error)
	DeleteAlbum(id int) error
}

// AlbumService implementa el servicio utilizando JSONPlaceholder
type AlbumService struct {
	restClient restclient.IRestClient
}

// GetAlbums obtiene albums desde JSONPlaceholder
func (s *AlbumService) GetAlbums(title string, userID int, id int) (*[]data.Album, error) {
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
	var albums []data.Album
	err = json.NewDecoder(resp.Body).Decode(&albums)
	if err != nil {
		return nil, err
	}
	return &albums, nil
}

func (s *AlbumService) GetAlbum(id int) (*data.Album, error) {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "GET"
	resp, err := s.restClient.NewRequest(method, url, bytes.NewBuffer(nil), nil)
	if err != nil {
		return nil, err
	}
	var album data.Album
	err = json.NewDecoder(resp.Body).Decode(&album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (s *AlbumService) CreateAlbum(album data.Album) (*data.Album, error) {
	validate := validator.New()
	if err := validate.Struct(album); err != nil {
		return nil, fmt.Errorf("Album can´t be created. Status Code: %s", err)
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
		return nil, fmt.Errorf("Album can´t be created. Status Code: %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Album can´t be created. Status Code: %d", resp.StatusCode)
	}
	var createdAlbum data.Album
	err = json.NewDecoder(resp.Body).Decode(&createdAlbum)
	if err != nil {
		return nil, err
	}
	return &createdAlbum, nil
}

func (s *AlbumService) UpdateAlbum(id int, album data.Album) (*data.Album, error) {
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
		return nil, fmt.Errorf("Album can´t be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedAlbum data.Album
	err = json.NewDecoder(resp.Body).Decode(&updatedAlbum)
	if err != nil {
		return nil, err
	}
	return &updatedAlbum, nil
}

func (s *AlbumService) PatchAlbum(id int, album data.Album) (*data.Album, error) {
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
		return nil, fmt.Errorf("Album can't be updated. Status Code: %d", resp.StatusCode)
	}
	var updatedAlbum data.Album
	err = json.NewDecoder(resp.Body).Decode(&updatedAlbum)
	if err != nil {
		return nil, err
	}
	return &updatedAlbum, nil
}

func (s *AlbumService) DeleteAlbum(id int) error {
	url := fmt.Sprintf("%s/%d", baseUrl, id)
	method := "DELETE"
	payload := bytes.NewBuffer(nil)
	resp, err := s.restClient.NewRequest(method, url, payload, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Album cannnt be created. Status Code: %d", resp.StatusCode)
	}
	return nil
}

// NewAlbumAlbumService crea una nueva instancia del servicio de albums
func NewAlbumService(client restclient.IRestClient) IAlbumService {
	return &AlbumService{
		restClient: client,
	}
}
