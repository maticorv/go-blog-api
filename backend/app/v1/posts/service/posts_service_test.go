package posts

import (
	"blog-api/app/mocks"
	data "blog-api/data"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockReader struct {
	*bytes.Reader
}

func (r *MockReader) Close() error {
	return nil
}

func TestPostService_GetPosts_Success(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&MockReader{bytes.NewReader([]byte(`[{"id": 1, "title": "Test Title", "userId": 10}]`))}),
	}
	//mockResp.Body.Read([]byte(`[{"id": 1, "title": "Test Title", "userId": 10}]`))
	mockClient := mocks.NewIRestClient(t)
	mockClient.On("NewRequest", "GET", "https://jsonplaceholder.typicode.com/posts?id=1&title=Test+Title&userId=10", mock.Anything, mock.Anything).Return(mockResp, nil)
	//
	service := &PostService{restClient: mockClient}
	posts, err := service.GetPosts("Test Title", 10, 1)
	require.NoError(t, err)
	require.NotNil(t, posts)
	require.Equal(t, 1, len(*posts))
}

func TestPostService_GetPosts_ErrorInRequest(t *testing.T) {
	mockClient := mocks.NewIRestClient(t)
	mockClient.On("NewRequest", "GET", "https://jsonplaceholder.typicode.com/posts?id=1&title=Test+Title&userId=10", mock.Anything, mock.Anything).Return(nil, errors.New("request error"))

	service := &PostService{restClient: mockClient}
	_, err := service.GetPosts("Test Title", 10, 1)
	require.Error(t, err)
}

func TestPostService_GetPost_Success(t *testing.T) {
	// Mock expected post data
	mockPost := &data.Post{
		ID:     1,
		Title:  "Test Title",
		UserID: 10,
	}
	// Mock client response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"id": 1, "title": "Test Title", "userId": 10}`))),
	}
	mockClient := mocks.NewIRestClient(t)
	mockClient.On("NewRequest", "GET", fmt.Sprintf("%s/%d", baseUrl, mockPost.ID), mock.Anything, mock.Anything).Return(mockResp, nil)
	// Create service and call GetPost
	service := &PostService{restClient: mockClient}
	post, err := service.GetPost(mockPost.ID)
	require.NoError(t, err)
	require.NotNil(t, post)
	require.Equal(t, *mockPost, *post)
}

func TestCreatePost(t *testing.T) {
	mockClient := mocks.NewIRestClient(t)
	service := &PostService{restClient: mockClient}
	post := data.Post{
		ID:     1,
		Title:  "Test Title",
		UserID: 10,
	}
	t.Run("success", func(t *testing.T) {
		mockClient.On("NewRequest", "POST", baseUrl, mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"id": 1, "title": "Test Title", "userId": 10}`))),
		}, nil)

		createdPost, err := service.CreatePost(post)
		assert.NoError(t, err)
		assert.Equal(t, post, *createdPost)
	})
	t.Run("error in validation", func(t *testing.T) {
		invalidPost := data.Post{
			ID:     1,
			Title:  "", // Title is required
			UserID: 10,
		}
		_, err := service.CreatePost(invalidPost)
		assert.Error(t, err)
	})
	t.Run("error in request", func(t *testing.T) {
		mockClient.On("NewRequest", "POST", baseUrl, mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(nil, errors.New("request error"))

		_, err := service.CreatePost(post)
		assert.Error(t, err)
	})
	t.Run("error in response", func(t *testing.T) {
		mockClient.On("NewRequest", "POST", baseUrl, mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)

		_, err := service.CreatePost(post)
		assert.Error(t, err)
	})
	t.Run("error in decoding", func(t *testing.T) {
		mockClient.On("NewRequest", "POST", baseUrl, mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`not a valid json`))),
		}, nil)
		_, err := service.CreatePost(post)
		assert.Error(t, err)
	})
}

func TestUpdatePost(t *testing.T) {
	mockClient := mocks.NewIRestClient(t)
	service := &PostService{restClient: mockClient}

	post := data.Post{
		ID:     1,
		Title:  "Updated Title",
		UserID: 10,
	}

	t.Run("success", func(t *testing.T) {
		mockClient.On("NewRequest", "PUT", fmt.Sprintf("%s/%d", baseUrl, post.ID), mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"id": 1, "title": "Updated Title", "userId": 10}`))),
		}, nil)

		updatedPost, err := service.UpdatePost(post.ID, post)
		assert.NoError(t, err)
		assert.Equal(t, post, *updatedPost)
	})

	t.Run("error in request", func(t *testing.T) {
		mockClient.On("NewRequest", "PUT", fmt.Sprintf("%s/%d", baseUrl, post.ID), mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(nil, errors.New("request error"))

		_, err := service.UpdatePost(post.ID, post)
		assert.Error(t, err)
	})

	t.Run("error in response", func(t *testing.T) {
		mockClient.On("NewRequest", "PUT", fmt.Sprintf("%s/%d", baseUrl, post.ID), mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)

		_, err := service.UpdatePost(post.ID, post)
		assert.Error(t, err)
	})

	t.Run("error in decoding", func(t *testing.T) {
		mockClient.On("NewRequest", "PUT", fmt.Sprintf("%s/%d", baseUrl, post.ID), mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`not a valid json`))),
		}, nil)

		_, err := service.UpdatePost(post.ID, post)
		assert.Error(t, err)
	})
}

func TestDeletePost(t *testing.T) {
	mockClient := mocks.NewIRestClient(t)
	service := &PostService{restClient: mockClient}
	postID := 1
	t.Run("success", func(t *testing.T) {
		mockClient.On("NewRequest", "DELETE", fmt.Sprintf("%s/%d", baseUrl, postID), mock.Anything, mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
		}, nil)
		err := service.DeletePost(postID)
		assert.NoError(t, err)
	})
}
