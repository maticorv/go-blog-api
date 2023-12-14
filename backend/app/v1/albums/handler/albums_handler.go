package albums

import (
	albums "blog-api/app/v1/albums/service"
	"blog-api/data"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// AlbumHandler maneja las solicitudes relacionadas con albums
type AlbumHandler struct {
	postService albums.IAlbumService
}

// NewAlbumHandler crea una nueva instancia del manejador de albums
func NewAlbumHandler(postService albums.IAlbumService) *AlbumHandler {
	return &AlbumHandler{
		postService: postService,
	}
}

// GetAlbums godoc
// @Description  Handler to get albums
// @Tags Albums
// @Description.markdown get albums
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       ///v1/albums [get] .
func (ph *AlbumHandler) GetAlbums(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userID, _ := strconv.Atoi(r.URL.Query().Get("userId"))
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	albums, err := ph.postService.GetAlbums(title, userID, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}

// GetAlbum godoc
// @Description Handler to get post
// @Tags Albums
// @Description.markdown get post
// @Accept		 json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      404
// @Router       ///v1/post/{postId} [get] .
func (ph *AlbumHandler) GetAlbum(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	post, err := ph.postService.GetAlbum(id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// CreateAlbum GetAlbum godoc
// @Description Handler to get post
// @Tags Albums
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
func (ph *AlbumHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var post data.Album
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {

		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	createdAlbum, err := ph.postService.CreateAlbum(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAlbum)
}

// UpdateAlbum godoc
// @Description Handler to update post
// @Tags Albums
// @Description.markdown update post
// @Accept		 json
// @UrlParam  		AlbumID path string true "AlbumID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [post] .
func (ph *AlbumHandler) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.Album
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	updatedAlbum, err := ph.postService.UpdateAlbum(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedAlbum)
}

// PatchAlbum godoc
// @Description Handler to patch post
// @Tags Albums
// @Description.markdown patch post
// @Accept		 json
// @UrlParam  		AlbumID path string true "AlbumID"
// @Param  		UserID path string true "UserID"
// @Param  		Title path string true "Title"
// @Param  		Body path string true "Body"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *AlbumHandler) PatchAlbum(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var post data.Album
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	patchedAlbum, err := ph.postService.UpdateAlbum(id, post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patchedAlbum)
}

// DeleteAlbum godoc
// @Description Handler to delete post
// @Tags Albums
// @Description.markdown delete post
// @Accept		 json
// @UrlParam  		AlbumID path string true "AlbumID"
// @Produce      json
// @Success      204
// @Failure      400
// @Failure      500
// @Router       ///v1/post/{postId} [patch] .
func (ph *AlbumHandler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = ph.postService.DeleteAlbum(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
