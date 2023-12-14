package comments

import (
	comments "blog-api/app/v1/comments/service"
	"blog-api/data"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// CommentHandler maneja las solicitudes relacionadas con comments
type CommentHandler struct {
	postService comments.ICommentService
}

// NewCommentHandler crea una nueva instancia del manejador de comments
func NewCommentHandler(postService comments.ICommentService) *CommentHandler {
	return &CommentHandler{
		postService: postService,
	}
}

// GetComments godoc
// @Description  Handler to get comments
// @Tags Comments
// @Description.markdown get comments
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       ///v1/comments [get] .
func (ph *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userID, _ := strconv.Atoi(r.URL.Query().Get("userId"))
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	comments, err := ph.postService.GetComments(title, userID, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// GetComment godoc
// @Description Handler to get post
// @Tags Comments
// @Description.markdown get post
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      404
// @Router       ///v1/post/{postId} [get] .
func (ph *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	post, err := ph.postService.GetComment(id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// CreateComment GetComment godoc
// @Description Handler to get post
// @Tags Comments
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
func (ph *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var post data.Comment
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {

		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	createdComment, err := ph.postService.CreateComment(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdComment)
}

// UpdateComment godoc
// @Description Handler to update post
// @Tags Comments
// @Description.markdown update post
// @Accept		 json
// @UrlParam  		CommentID path string true "CommentID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [post] .
func (ph *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.Comment
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	updatedComment, err := ph.postService.UpdateComment(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedComment)
}

// PatchComment godoc
// @Description Handler to patch post
// @Tags Comments
// @Description.markdown patch post
// @Accept		 json
// @UrlParam  		CommentID path string true "CommentID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *CommentHandler) PatchComment(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.Comment
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	patchedComment, err := ph.postService.UpdateComment(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patchedComment)
}

// DeleteComment godoc
// @Description Handler to delete post
// @Tags Comments
// @Description.markdown delete post
// @Accept		 json
// @UrlParam  		CommentID path string true "CommentID"
// @Produce      json
// @Success      204
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = ph.postService.DeleteComment(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
