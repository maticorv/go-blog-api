package users

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

func TestGetUsers_Success(t *testing.T) {
	// Mock post service
	mockUserService := mocks.NewIUserService(t)
	defer mockUserService.AssertExpectations(t)
	// Mock expected users
	mockUsers := []data.User{{ID: 1, Name: "My User", Username: "myuser", Email: "asd@gmail.com", Address: data.Address{}, Phone: "", Website: "", Company: data.Company{}}}
	// Set mock expectations
	mockUserService.On("GetUsers", "", 0, 0).Return(&mockUsers, nil)
	// Create handler and request
	handler := UserHandler{postService: mockUserService}
	req, _ := http.NewRequest("GET", "/users", nil)
	// Create mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetUsers(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockUsers)
	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestGetUser_Success(t *testing.T) {
	// Mock post service
	ctx := chi.NewRouteContext()
	mockUserService := mocks.NewIUserService(t)
	defer mockUserService.AssertExpectations(t)
	// Mock expected post
	mockUser := data.User{ID: 1, Name: "My User", Username: "myuser", Email: "asd@gmail.com", Address: data.Address{}, Phone: "", Website: "", Company: data.Company{}}
	// Set mock expectations
	mockUserService.On("GetUser", 1).Return(&mockUser, nil)
	// Create handler and request
	handler := &UserHandler{postService: mockUserService}
	req, _ := http.NewRequest("GET", "/users/1", nil)
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.GetUser(mockRecorder, req)
	jsonBytes, _ := json.Marshal(mockUser)

	// Assertions
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestCreateUser_Success(t *testing.T) {
	// Mock post service
	mockUserService := mocks.NewIUserService(t)
	defer mockUserService.AssertExpectations(t)
	// Mock expected post
	mockUser := data.User{ID: 1, Name: "My User", Username: "myuser", Email: "asd@gmail.com", Address: data.Address{}, Phone: "", Website: "", Company: data.Company{}}
	mockCreatedUser := data.User{ID: 1, Name: "My User", Username: "myuser", Email: "asd@gmail.com", Address: data.Address{}, Phone: "", Website: "", Company: data.Company{}}

	// Set mock expectations
	mockUserService.On("CreateUser", mockUser).Return(&mockCreatedUser, nil)
	// Create handler and request
	handler := &UserHandler{postService: mockUserService}
	reqBody, err := json.Marshal(mockUser)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
	// Mock response recorder
	mockRecorder := httptest.NewRecorder()
	// Handle request
	handler.CreateUser(mockRecorder, req)
	// Assertions
	jsonBytes, _ := json.Marshal(mockCreatedUser)
	require.Equal(t, http.StatusCreated, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestUpdateUser_Success(t *testing.T) {
	// Mock post service
	mockUserService := mocks.NewIUserService(t)
	defer mockUserService.AssertExpectations(t)

	// Mock expected post
	mockUser := data.User{ID: 1, Name: "My User", Username: "myuser", Email: "asd@gmail.com", Address: data.Address{}, Phone: "", Website: "", Company: data.Company{}}
	mockUpdatedUser := data.User{ID: 1, Name: "My User", Username: "myuser", Email: "asd@gmail.com", Address: data.Address{}, Phone: "", Website: "", Company: data.Company{}}

	// Set mock expectations
	mockUserService.On("UpdateUser", 1, mockUser).Return(&mockUpdatedUser, nil)

	// Create handler and request
	handler := &UserHandler{postService: mockUserService}
	reqBody, err := json.Marshal(mockUser)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.UpdateUser(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockUpdatedUser)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestPatchUser_Success(t *testing.T) {
	// Mock post service
	mockUserService := mocks.NewIUserService(t)
	defer mockUserService.AssertExpectations(t)

	// Mock expected post with only updated title
	mockUser := data.User{ID: 1, Name: "My User", Username: "myuser", Email: "asd@gmail.com", Address: data.Address{}, Phone: "", Website: "", Company: data.Company{}}
	mockPatchedUser := data.User{ID: 1, Name: "My User", Username: "myuser", Email: "asd@gmail.com", Address: data.Address{}, Phone: "", Website: "", Company: data.Company{}}

	// Set mock expectations
	mockUserService.On("UpdateUser", 1, mockUser).Return(&mockPatchedUser, nil)

	// Create handler and request with only title in the body
	handler := &UserHandler{postService: mockUserService}
	reqBody, err := json.Marshal(mockUser)
	require.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader(reqBody))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.PatchUser(mockRecorder, req)

	// Assertions
	jsonBytes, _ := json.Marshal(mockPatchedUser)
	require.Equal(t, http.StatusOK, mockRecorder.Code)
	require.Equal(t, "application/json", mockRecorder.HeaderMap["Content-Type"][0])
	require.Equal(t, string(jsonBytes)+"\n", mockRecorder.Body.String())
}

func TestDeleteUser_Success(t *testing.T) {
	// Mock post service
	mockUserService := mocks.NewIUserService(t)
	defer mockUserService.AssertExpectations(t)

	// Set mock expectations
	mockUserService.On("DeleteUser", 1).Return(nil)

	// Create handler and request
	handler := &UserHandler{postService: mockUserService}
	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// Mock response recorder
	mockRecorder := httptest.NewRecorder()

	// Handle request
	handler.DeleteUser(mockRecorder, req)

	// Assertions
	require.Equal(t, http.StatusNoContent, mockRecorder.Code)
}
