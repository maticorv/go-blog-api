package restclient

import (
	"io"
	"net/http"
)

type IRestClient interface {
	NewRequest(method string, url string, body io.Reader, headers map[string]string) (*http.Response, error)
}
type RestClient struct{}

// NewRequest RestClient realiza una solicitud HTTP con los parámetros dados
func (rc *RestClient) NewRequest(method string, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()
	return resp, nil
}

func NewRestClient() *RestClient {
	return &RestClient{}
}
