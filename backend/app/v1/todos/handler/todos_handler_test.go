package todos

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

func TestGetTodos_Success(t *testing.T) {
	// Mock post service
	mockTodoService := mocks.NewITodoService(t)
	defer mockTodoService.AssertExpectations(t)
	// Mock expected todos
	mockTodos := []data.Todo{{ID: 1, Title: "My Todo", UserID: 1}}
	// Set mock expectations
	mockTodoService.On("GetTodos", "", 0, 0).Return(&mockTodos, nil)
	// Create handler and request
	handler := TodoHandler{postService: mockTodoService}
	req, _ := http.NewRequest("GET", "/todos", nil)
	// Create mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetTodos(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockTodos)
	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestGetTodo_Success(t *testing.T) {
	// Mock post service
	ctx := chi.NewRouteContext()
	mockTodoService := mocks.NewITodoService(t)
	defer mockTodoService.AssertExpectations(t)
	// Mock expected post
	mockTodo := data.Todo{ID: 1, Title: "My Todo", UserID: 1}
	// Set mock expectations
	mockTodoService.On("GetTodo", 1).Return(&mockTodo, nil)
	// Create handler and request
	handler := &TodoHandler{postService: mockTodoService}
	req, _ := http.NewRequest("GET", "/todos/1", nil)
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetTodo(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockTodo)

	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestCreateTodo_Success(t *testing.T) {
	// Mock post service
	mockTodoService := mocks.NewITodoService(t)
	defer mockTodoService.AssertExpectations(t)
	// Mock expected post
	mockTodo := data.Todo{Title: "My Todo", UserID: 1}
	mockCreatedTodo := data.Todo{ID: 1, Title: "My Todo", UserID: 1}
	// Set mock expectations
	mockTodoService.On("CreateTodo", mockTodo).Return(&mockCreatedTodo, nil)
	// Create handler and request
	handler := &TodoHandler{postService: mockTodoService}
	reqBody, err := json.Marshal(mockTodo)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewReader(reqBody))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.CreateTodo(mockRecorder, req)
	// Assertions
	jsonBytes, _ := json.Marshal(mockCreatedTodo)
	require.Equal(t, http.StatusCreated, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestUpdateTodo_Success(t *testing.T) {
	// Mock post service
	mockTodoService := mocks.NewITodoService(t)
	defer mockTodoService.AssertExpectations(t)

	// Mock expected post
	mockTodo := data.Todo{ID: 1, Title: "Updated Todo", UserID: 1}
	mockUpdatedTodo := data.Todo{ID: 1, Title: "Updated Todo", UserID: 1}

	// Set mock expectations
	mockTodoService.On("UpdateTodo", 1, mockTodo).Return(&mockUpdatedTodo, nil)

	// Create handler and request
	handler := &TodoHandler{postService: mockTodoService}
	reqBody, err := json.Marshal(mockTodo)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.UpdateTodo(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockUpdatedTodo)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestPatchTodo_Success(t *testing.T) {
	// Mock post service
	mockTodoService := mocks.NewITodoService(t)
	defer mockTodoService.AssertExpectations(t)

	// Mock expected post with only updated title
	mockTodo := data.Todo{Title: "Updated Title"}
	mockPatchedTodo := data.Todo{ID: 1, Title: "Updated Title", UserID: 1}

	// Set mock expectations
	mockTodoService.On("UpdateTodo", 1, mockTodo).Return(&mockPatchedTodo, nil)

	// Create handler and request with only title in the body
	handler := &TodoHandler{postService: mockTodoService}
	reqBody, err := json.Marshal(mockTodo)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPatch, "/todos/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.PatchTodo(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockPatchedTodo)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestDeleteTodo_Success(t *testing.T) {
	// Mock post service
	mockTodoService := mocks.NewITodoService(t)
	defer mockTodoService.AssertExpectations(t)

	// Set mock expectations
	mockTodoService.On("DeleteTodo", 1).Return(nil)

	// Create handler and request
	handler := &TodoHandler{postService: mockTodoService}
	req, _ := http.NewRequest(http.MethodDelete, "/todos/1", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.DeleteTodo(mockRecorder, req)

	// Assertions
	require.Equal(t, http.StatusNoContent, mockRecorder.Code)
}
