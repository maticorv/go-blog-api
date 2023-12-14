package users

import (
	users "blog-api/app/v1/users/service"
	"blog-api/data"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// UserHandler maneja las solicitudes relacionadas con users
type UserHandler struct {
	postService users.IUserService
}

// NewUserHandler crea una nueva instancia del manejador de users
func NewUserHandler(postService users.IUserService) *UserHandler {
	return &UserHandler{
		postService: postService,
	}
}

// GetUsers godoc
// @Description  Handler to get users
// @Tags Users
// @Description.markdown get users
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       ///v1/users [get] .
func (ph *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userID, _ := strconv.Atoi(r.URL.Query().Get("userId"))
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	users, err := ph.postService.GetUsers(title, userID, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser godoc
// @Description Handler to get post
// @Tags Users
// @Description.markdown get post
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      404
// @Router       ///v1/post/{postId} [get] .
func (ph *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	post, err := ph.postService.GetUser(id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// CreateUser GetUser godoc
// @Description Handler to get post
// @Tags Users
// @Description.markdown creat post
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Accept		 json
// @Produce      json
// @Success      201
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [post] .
func (ph *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var post data.User
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {

		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	createdUser, err := ph.postService.CreateUser(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// UpdateUser godoc
// @Description Handler to update post
// @Tags Users
// @Description.markdown update post
// @Accept		 json
// @UrlParam  		UserID path string true "UserID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [post] .
func (ph *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.User
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	updatedUser, err := ph.postService.UpdateUser(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// PatchUser godoc
// @Description Handler to patch post
// @Tags Users
// @Description.markdown patch post
// @Accept		 json
// @UrlParam  		UserID path string true "UserID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.User
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	patchedUser, err := ph.postService.UpdateUser(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patchedUser)
}

// DeleteUser godoc
// @Description Handler to delete post
// @Tags Users
// @Description.markdown delete post
// @Accept		 json
// @UrlParam  		UserID path string true "UserID"
// @Produce      json
// @Success      204
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = ph.postService.DeleteUser(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
