package albums

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

func TestGetAlbums_Success(t *testing.T) {
	// Mock post service
	mockAlbumService := mocks.NewIAlbumService(t)
	defer mockAlbumService.AssertExpectations(t)
	// Mock expected albums
	mockAlbums := []data.Album{{ID: 1, Title: "My Album", UserID: 1}}
	// Set mock expectations
	mockAlbumService.On("GetAlbums", "", 0, 0).Return(&mockAlbums, nil)
	// Create handler and request
	handler := AlbumHandler{postService: mockAlbumService}
	req, _ := http.NewRequest("GET", "/albums", nil)
	// Create mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetAlbums(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockAlbums)
	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestGetAlbum_Success(t *testing.T) {
	// Mock post service
	ctx := chi.NewRouteContext()
	mockAlbumService := mocks.NewIAlbumService(t)
	defer mockAlbumService.AssertExpectations(t)
	// Mock expected post
	mockAlbum := data.Album{ID: 1, Title: "My Album", UserID: 1}
	// Set mock expectations
	mockAlbumService.On("GetAlbum", 1).Return(&mockAlbum, nil)
	// Create handler and request
	handler := &AlbumHandler{postService: mockAlbumService}
	req, _ := http.NewRequest("GET", "/albums/1", nil)
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetAlbum(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockAlbum)

	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestCreateAlbum_Success(t *testing.T) {
	// Mock post service
	mockAlbumService := mocks.NewIAlbumService(t)
	defer mockAlbumService.AssertExpectations(t)
	// Mock expected post
	mockAlbum := data.Album{Title: "My Album", UserID: 1}
	mockCreatedAlbum := data.Album{ID: 1, Title: "My Album", UserID: 1}
	// Set mock expectations
	mockAlbumService.On("CreateAlbum", mockAlbum).Return(&mockCreatedAlbum, nil)
	// Create handler and request
	handler := &AlbumHandler{postService: mockAlbumService}
	reqBody, err := json.Marshal(mockAlbum)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/albums", bytes.NewReader(reqBody))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.CreateAlbum(mockRecorder, req)
	// Assertions
	jsonBytes, _ := json.Marshal(mockCreatedAlbum)
	require.Equal(t, http.StatusCreated, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestUpdateAlbum_Success(t *testing.T) {
	// Mock post service
	mockAlbumService := mocks.NewIAlbumService(t)
	defer mockAlbumService.AssertExpectations(t)

	// Mock expected post
	mockAlbum := data.Album{ID: 1, Title: "Updated Album", UserID: 1}
	mockUpdatedAlbum := data.Album{ID: 1, Title: "Updated Album", UserID: 1}

	// Set mock expectations
	mockAlbumService.On("UpdateAlbum", 1, mockAlbum).Return(&mockUpdatedAlbum, nil)

	// Create handler and request
	handler := &AlbumHandler{postService: mockAlbumService}
	reqBody, err := json.Marshal(mockAlbum)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPut, "/albums/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.UpdateAlbum(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockUpdatedAlbum)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestPatchAlbum_Success(t *testing.T) {
	// Mock post service
	mockAlbumService := mocks.NewIAlbumService(t)
	defer mockAlbumService.AssertExpectations(t)

	// Mock expected post with only updated title
	mockAlbum := data.Album{Title: "Updated Title"}
	mockPatchedAlbum := data.Album{ID: 1, Title: "Updated Title", UserID: 1}

	// Set mock expectations
	mockAlbumService.On("UpdateAlbum", 1, mockAlbum).Return(&mockPatchedAlbum, nil)

	// Create handler and request with only title in the body
	handler := &AlbumHandler{postService: mockAlbumService}
	reqBody, err := json.Marshal(mockAlbum)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPatch, "/albums/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.PatchAlbum(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockPatchedAlbum)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestDeleteAlbum_Success(t *testing.T) {
	// Mock post service
	mockAlbumService := mocks.NewIAlbumService(t)
	defer mockAlbumService.AssertExpectations(t)

	// Set mock expectations
	mockAlbumService.On("DeleteAlbum", 1).Return(nil)

	// Create handler and request
	handler := &AlbumHandler{postService: mockAlbumService}
	req, _ := http.NewRequest(http.MethodDelete, "/albums/1", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.DeleteAlbum(mockRecorder, req)

	// Assertions
	require.Equal(t, http.StatusNoContent, mockRecorder.Code)
}
