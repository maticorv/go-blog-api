package posts

import (
	posts "blog-api/app/v1/posts/service"
	"blog-api/data"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// PostHandler maneja las solicitudes relacionadas con posts
type PostHandler struct {
	postService posts.IPostService
}

// NewPostHandler crea una nueva instancia del manejador de posts
func NewPostHandler(postService posts.IPostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// GetPosts godoc
// @Description  Handler to get posts
// @Tags Posts
// @Description.markdown get posts
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       ///v1/posts [get] .
func (ph *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userID, _ := strconv.Atoi(r.URL.Query().Get("userId"))
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	posts, err := ph.postService.GetPosts(title, userID, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetPost godoc
// @Description Handler to get post
// @Tags Posts
// @Description.markdown get post
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      404
// @Router       ///v1/post/{postId} [get] .
func (ph *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	post, err := ph.postService.GetPost(id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// CreatePost GetPost godoc
// @Description Handler to get post
// @Tags Posts
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
func (ph *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post data.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {

		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	createdPost, err := ph.postService.CreatePost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPost)
}

// UpdatePost godoc
// @Description Handler to update post
// @Tags Posts
// @Description.markdown update post
// @Accept		 json
// @UrlParam  		PostID path string true "PostID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [post] .
func (ph *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	updatedPost, err := ph.postService.UpdatePost(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPost)
}

// PatchPost godoc
// @Description Handler to patch post
// @Tags Posts
// @Description.markdown patch post
// @Accept		 json
// @UrlParam  		PostID path string true "PostID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *PostHandler) PatchPost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	patchedPost, err := ph.postService.UpdatePost(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patchedPost)
}

// DeletePost godoc
// @Description Handler to delete post
// @Tags Posts
// @Description.markdown delete post
// @Accept		 json
// @UrlParam  		PostID path string true "PostID"
// @Produce      json
// @Success      204
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = ph.postService.DeletePost(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
