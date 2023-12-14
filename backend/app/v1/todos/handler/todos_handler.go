package todos

import (
	todos "blog-api/app/v1/todos/service"
	"blog-api/data"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// TodoHandler maneja las solicitudes relacionadas con todos
type TodoHandler struct {
	postService todos.ITodoService
}

// NewTodoHandler crea una nueva instancia del manejador de todos
func NewTodoHandler(postService todos.ITodoService) *TodoHandler {
	return &TodoHandler{
		postService: postService,
	}
}

// GetTodos godoc
// @Description  Handler to get todoss
// @Tags Todos
// @Description.markdown get todos
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       ///v1/todos [get] .
func (ph *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userID, _ := strconv.Atoi(r.URL.Query().Get("userId"))
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	todos, err := ph.postService.GetTodos(title, userID, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// GetTodo godoc
// @Description Handler to get post
// @Tags Todos
// @Description.markdown get post
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      404
// @Router       ///v1/post/{postId} [get] .
func (ph *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	post, err := ph.postService.GetTodo(id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// CreateTodo GetTodo godoc
// @Description Handler to get post
// @Tags Todos
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
func (ph *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var post data.Todo
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {

		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	createdTodo, err := ph.postService.CreateTodo(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTodo)
}

// UpdateTodo godoc
// @Description Handler to update post
// @Tags Todos
// @Description.markdown update post
// @Accept		 json
// @UrlParam  		TodoID path string true "TodoID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [post] .
func (ph *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.Todo
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	updatedTodo, err := ph.postService.UpdateTodo(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

// PatchTodo godoc
// @Description Handler to patch post
// @Tags Todos
// @Description.markdown patch post
// @Accept		 json
// @UrlParam  		TodoID path string true "TodoID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *TodoHandler) PatchTodo(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.Todo
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	patchedTodo, err := ph.postService.UpdateTodo(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patchedTodo)
}

// DeleteTodo godoc
// @Description Handler to delete post
// @Tags Todos
// @Description.markdown delete post
// @Accept		 json
// @UrlParam  		TodoID path string true "TodoID"
// @Produce      json
// @Success      204
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = ph.postService.DeleteTodo(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
