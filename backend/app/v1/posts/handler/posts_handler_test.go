package posts

import (
	"blog-api/app/mocks"
	"blog-api/data"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPosts_Success(t *testing.T) {
	// Mock post service
	mockPostService := mocks.NewIPostService(t)
	defer mockPostService.AssertExpectations(t)
	// Mock expected posts
	mockPosts := []data.Post{{ID: 1, Title: "My Post", Body: "This is my post", UserID: 1}}
	// Set mock expectations
	mockPostService.On("GetPosts", "", 0, 0).Return(&mockPosts, nil)
	// Create handler and request
	handler := PostHandler{postService: mockPostService}
	req, _ := http.NewRequest("GET", "/posts", nil)
	// Create mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetPosts(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockPosts)
	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestGetPost_Success(t *testing.T) {
	// Mock post service
	ctx := chi.NewRouteContext()
	mockPostService := mocks.NewIPostService(t)
	defer mockPostService.AssertExpectations(t)
	// Mock expected post
	mockPost := data.Post{ID: 1, Title: "My Post", Body: "This is my post", UserID: 1}
	// Set mock expectations
	mockPostService.On("GetPost", 1).Return(&mockPost, nil)
	// Create handler and request
	handler := &PostHandler{postService: mockPostService}
	req, _ := http.NewRequest("GET", "/posts/1", nil)
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetPost(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockPost)

	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestCreatePost_Success(t *testing.T) {
	// Mock post service
	mockPostService := mocks.NewIPostService(t)
	defer mockPostService.AssertExpectations(t)
	// Mock expected post
	mockPost := data.Post{Title: "My Post", Body: "This is my post", UserID: 1}
	mockCreatedPost := data.Post{ID: 1, Title: "My Post", Body: "This is my post", UserID: 1}
	// Set mock expectations
	mockPostService.On("CreatePost", mockPost).Return(&mockCreatedPost, nil)
	// Create handler and request
	handler := &PostHandler{postService: mockPostService}
	reqBody, err := json.Marshal(mockPost)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewReader(reqBody))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.CreatePost(mockRecorder, req)
	// Assertions
	jsonBytes, _ := json.Marshal(mockCreatedPost)
	require.Equal(t, http.StatusCreated, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestUpdatePost_Success(t *testing.T) {
	// Mock post service
	mockPostService := mocks.NewIPostService(t)
	defer mockPostService.AssertExpectations(t)

	// Mock expected post
	mockPost := data.Post{ID: 1, Title: "Updated Post", Body: "This is an updated post", UserID: 1}
	mockUpdatedPost := data.Post{ID: 1, Title: "Updated Post", Body: "This is an updated post", UserID: 1}

	// Set mock expectations
	mockPostService.On("UpdatePost", 1, mockPost).Return(&mockUpdatedPost, nil)

	// Create handler and request
	handler := &PostHandler{postService: mockPostService}
	reqBody, err := json.Marshal(mockPost)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPut, "/posts/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.UpdatePost(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockUpdatedPost)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestPatchPost_Success(t *testing.T) {
	// Mock post service
	mockPostService := mocks.NewIPostService(t)
	defer mockPostService.AssertExpectations(t)

	// Mock expected post with only updated title
	mockPost := data.Post{Title: "Updated Title"}
	mockPatchedPost := data.Post{ID: 1, Title: "Updated Title", Body: "This is my post", UserID: 1}

	// Set mock expectations
	mockPostService.On("UpdatePost", 1, mockPost).Return(&mockPatchedPost, nil)

	// Create handler and request with only title in the body
	handler := &PostHandler{postService: mockPostService}
	reqBody, err := json.Marshal(mockPost)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPatch, "/posts/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.PatchPost(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockPatchedPost)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestDeletePost_Success(t *testing.T) {
	// Mock post service
	mockPostService := mocks.NewIPostService(t)
	defer mockPostService.AssertExpectations(t)

	// Set mock expectations
	mockPostService.On("DeletePost", 1).Return(nil)

	// Create handler and request
	handler := &PostHandler{postService: mockPostService}
	req, _ := http.NewRequest(http.MethodDelete, "/posts/1", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.DeletePost(mockRecorder, req)

	// Assertions
	require.Equal(t, http.StatusNoContent, mockRecorder.Code)
}
