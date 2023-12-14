package restclient

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRestClient_NewRequest_Success(t *testing.T) {
	// Mock body reader
	mockBody := bytes.NewBufferString("hello world")
	// Mock headers
	headers := map[string]string{}
	// Create client and request
	rc := &RestClient{}
	req, err := rc.NewRequest(http.MethodGet, "https://jsonplaceholder.typicode.com/posts", mockBody, headers)
	require.NoError(t, err)
	defer req.Body.Close()
	// Verify HTTP method and URL
	require.Equal(t, http.MethodGet, req.Request.Method)
	require.Equal(t, "https://jsonplaceholder.typicode.com/posts", req.Request.URL.String())

	// Verify headers
	for key, value := range headers {
		require.Equal(t, value, req.Header.Get(key))
	}

	// Verify response status code
	require.Equal(t, http.StatusOK, req.StatusCode)
}
