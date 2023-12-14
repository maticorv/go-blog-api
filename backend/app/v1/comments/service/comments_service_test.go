package comments

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

func TestCommentService_GetComments_Success(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&MockReader{bytes.NewReader([]byte(`[{"id": 1, "title": "Test Title", "userId": 10}]`))}),
	}
	//mockResp.Body.Read([]byte(`[{"id": 1, "title": "Test Title", "userId": 10}]`))
	mockClient := mocks.NewIRestClient(t)
	mockClient.On("NewRequest", "GET", "https://jsonplaceholder.typicode.com/comments?id=1&title=Test+Title&userId=10", mock.Anything, mock.Anything).Return(mockResp, nil)
	//
	service := &CommentService{restClient: mockClient}
	comments, err := service.GetComments("Test Title", 10, 1)
	require.NoError(t, err)
	require.NotNil(t, comments)
	require.Equal(t, 1, len(*comments))
}

func TestCommentService_GetComments_ErrorInRequest(t *testing.T) {
	mockClient := mocks.NewIRestClient(t)
	mockClient.On("NewRequest", "GET", "https://jsonplaceholder.typicode.com/comments?id=1&title=Test+Title&userId=10", mock.Anything, mock.Anything).Return(nil, errors.New("request error"))

	service := &CommentService{restClient: mockClient}
	_, err := service.GetComments("Test Title", 10, 1)
	require.Error(t, err)
}

func TestCommentService_GetComment_Success(t *testing.T) {
	// Mock expected album data
	mockComment := &data.Comment{
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
	mockClient.On("NewRequest", "GET", fmt.Sprintf("%s/%d", baseUrl, mockComment.ID), mock.Anything, mock.Anything).Return(mockResp, nil)
	// Create service and call GetComment
	service := &CommentService{restClient: mockClient}
	album, err := service.GetComment(mockComment.ID)
	require.NoError(t, err)
	require.NotNil(t, album)
	require.Equal(t, *mockComment, *album)
}

func TestCreateComment(t *testing.T) {
	mockClient := mocks.NewIRestClient(t)
	service := &CommentService{restClient: mockClient}
	album := data.Comment{
		ID:     1,
		Title:  "Test Title",
		UserID: 10,
	}
	t.Run("success", func(t *testing.T) {
		mockClient.On("NewRequest", "POST", baseUrl, mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"id": 1, "title": "Test Title", "userId": 10}`))),
		}, nil)

		createdComment, err := service.CreateComment(album)
		assert.NoError(t, err)
		assert.Equal(t, album, *createdComment)
	})
	t.Run("error in validation", func(t *testing.T) {
		invalidComment := data.Comment{
			ID:     1,
			Title:  "", // Title is required
			UserID: 10,
		}
		_, err := service.CreateComment(invalidComment)
		assert.Error(t, err)
	})
	t.Run("error in request", func(t *testing.T) {
		mockClient.On("NewRequest", "POST", baseUrl, mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(nil, errors.New("request error"))

		_, err := service.CreateComment(album)
		assert.Error(t, err)
	})
	t.Run("error in response", func(t *testing.T) {
		mockClient.On("NewRequest", "POST", baseUrl, mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)

		_, err := service.CreateComment(album)
		assert.Error(t, err)
	})
	t.Run("error in decoding", func(t *testing.T) {
		mockClient.On("NewRequest", "POST", baseUrl, mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`not a valid json`))),
		}, nil)
		_, err := service.CreateComment(album)
		assert.Error(t, err)
	})
}

func TestUpdateComment(t *testing.T) {
	mockClient := mocks.NewIRestClient(t)
	service := &CommentService{restClient: mockClient}

	album := data.Comment{
		ID:     1,
		Title:  "Updated Title",
		UserID: 10,
	}

	t.Run("success", func(t *testing.T) {
		mockClient.On("NewRequest", "PUT", fmt.Sprintf("%s/%d", baseUrl, album.ID), mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"id": 1, "title": "Updated Title", "userId": 10}`))),
		}, nil)

		updatedComment, err := service.UpdateComment(album.ID, album)
		assert.NoError(t, err)
		assert.Equal(t, album, *updatedComment)
	})

	t.Run("error in request", func(t *testing.T) {
		mockClient.On("NewRequest", "PUT", fmt.Sprintf("%s/%d", baseUrl, album.ID), mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(nil, errors.New("request error"))

		_, err := service.UpdateComment(album.ID, album)
		assert.Error(t, err)
	})

	t.Run("error in response", func(t *testing.T) {
		mockClient.On("NewRequest", "PUT", fmt.Sprintf("%s/%d", baseUrl, album.ID), mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)

		_, err := service.UpdateComment(album.ID, album)
		assert.Error(t, err)
	})

	t.Run("error in decoding", func(t *testing.T) {
		mockClient.On("NewRequest", "PUT", fmt.Sprintf("%s/%d", baseUrl, album.ID), mock.Anything, map[string]string{"Content-Type": " application/json"}).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`not a valid json`))),
		}, nil)

		_, err := service.UpdateComment(album.ID, album)
		assert.Error(t, err)
	})
}

func TestDeleteComment(t *testing.T) {
	mockClient := mocks.NewIRestClient(t)
	service := &CommentService{restClient: mockClient}
	albumID := 1
	t.Run("success", func(t *testing.T) {
		mockClient.On("NewRequest", "DELETE", fmt.Sprintf("%s/%d", baseUrl, albumID), mock.Anything, mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
		}, nil)
		err := service.DeleteComment(albumID)
		assert.NoError(t, err)
	})
}
