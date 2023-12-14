package comments

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

func TestGetComments_Success(t *testing.T) {
	// Mock post service
	mockCommentService := mocks.NewICommentService(t)
	defer mockCommentService.AssertExpectations(t)
	// Mock expected comments
	mockComments := []data.Comment{{ID: 1, Title: "My Comment", UserID: 1}}
	// Set mock expectations
	mockCommentService.On("GetComments", "", 0, 0).Return(&mockComments, nil)
	// Create handler and request
	handler := CommentHandler{postService: mockCommentService}
	req, _ := http.NewRequest("GET", "/comments", nil)
	// Create mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetComments(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockComments)
	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestGetComment_Success(t *testing.T) {
	// Mock post service
	ctx := chi.NewRouteContext()
	mockCommentService := mocks.NewICommentService(t)
	defer mockCommentService.AssertExpectations(t)
	// Mock expected post
	mockComment := data.Comment{ID: 1, Title: "My Comment", UserID: 1}
	// Set mock expectations
	mockCommentService.On("GetComment", 1).Return(&mockComment, nil)
	// Create handler and request
	handler := &CommentHandler{postService: mockCommentService}
	req, _ := http.NewRequest("GET", "/comments/1", nil)
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetComment(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockComment)

	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestCreateComment_Success(t *testing.T) {
	// Mock post service
	mockCommentService := mocks.NewICommentService(t)
	defer mockCommentService.AssertExpectations(t)
	// Mock expected post
	mockComment := data.Comment{Title: "My Comment", UserID: 1}
	mockCreatedComment := data.Comment{ID: 1, Title: "My Comment", UserID: 1}
	// Set mock expectations
	mockCommentService.On("CreateComment", mockComment).Return(&mockCreatedComment, nil)
	// Create handler and request
	handler := &CommentHandler{postService: mockCommentService}
	reqBody, err := json.Marshal(mockComment)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/comments", bytes.NewReader(reqBody))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.CreateComment(mockRecorder, req)
	// Assertions
	jsonBytes, _ := json.Marshal(mockCreatedComment)
	require.Equal(t, http.StatusCreated, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestUpdateComment_Success(t *testing.T) {
	// Mock post service
	mockCommentService := mocks.NewICommentService(t)
	defer mockCommentService.AssertExpectations(t)

	// Mock expected post
	mockComment := data.Comment{ID: 1, Title: "Updated Comment", UserID: 1}
	mockUpdatedComment := data.Comment{ID: 1, Title: "Updated Comment", UserID: 1}

	// Set mock expectations
	mockCommentService.On("UpdateComment", 1, mockComment).Return(&mockUpdatedComment, nil)

	// Create handler and request
	handler := &CommentHandler{postService: mockCommentService}
	reqBody, err := json.Marshal(mockComment)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPut, "/comments/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.UpdateComment(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockUpdatedComment)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestPatchComment_Success(t *testing.T) {
	// Mock post service
	mockCommentService := mocks.NewICommentService(t)
	defer mockCommentService.AssertExpectations(t)

	// Mock expected post with only updated title
	mockComment := data.Comment{Title: "Updated Title"}
	mockPatchedComment := data.Comment{ID: 1, Title: "Updated Title", UserID: 1}

	// Set mock expectations
	mockCommentService.On("UpdateComment", 1, mockComment).Return(&mockPatchedComment, nil)

	// Create handler and request with only title in the body
	handler := &CommentHandler{postService: mockCommentService}
	reqBody, err := json.Marshal(mockComment)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPatch, "/comments/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.PatchComment(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockPatchedComment)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestDeleteComment_Success(t *testing.T) {
	// Mock post service
	mockCommentService := mocks.NewICommentService(t)
	defer mockCommentService.AssertExpectations(t)

	// Set mock expectations
	mockCommentService.On("DeleteComment", 1).Return(nil)

	// Create handler and request
	handler := &CommentHandler{postService: mockCommentService}
	req, _ := http.NewRequest(http.MethodDelete, "/comments/1", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.DeleteComment(mockRecorder, req)

	// Assertions
	require.Equal(t, http.StatusNoContent, mockRecorder.Code)
}
