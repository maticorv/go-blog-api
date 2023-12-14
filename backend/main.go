// This example demonstrates the use of the render subpackage, with
// a quick concept for how to support multiple api versions.
package main

import (
	"blog-api/app/clients/restclient"
	ah "blog-api/app/v1/albums/handler"
	as "blog-api/app/v1/albums/service"
	ch "blog-api/app/v1/comments/handler"
	cs "blog-api/app/v1/comments/service"
	ph "blog-api/app/v1/posts/handler"
	ps "blog-api/app/v1/posts/service"
	th "blog-api/app/v1/todos/handler"
	ts "blog-api/app/v1/todos/service"
	uh "blog-api/app/v1/users/handler"
	us "blog-api/app/v1/users/service"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

// @title bolg-api example
// @version 1.0
// @description bolg-api
// @BasePath /
func main() {
	r := chi.NewRouter()
	restClient := restclient.NewRestClient()
	postService := ps.NewPostService(restClient)
	postHandler := ph.NewPostHandler(postService)
	albumsService := as.NewAlbumService(restClient)
	albumHandler := ah.NewAlbumHandler(albumsService)
	commentsService := cs.NewCommentService(restClient)
	commentHandler := ch.NewCommentHandler(commentsService)
	todoService := ts.NewTodoService(restClient)
	todoHandler := th.NewTodoHandler(todoService)
	userService := us.NewUserService(restClient)
	userHandler := uh.NewUserHandler(userService)
	// Logger
	logger := httplog.NewLogger("blog-api", httplog.Options{
		LogLevel: slog.LevelDebug,
		// JSON:             true,
		Concise: true,
		// RequestHeaders:   true,
		// ResponseHeaders:  true,
		MessageFieldName: "message",
		LevelFieldName:   "severity",
		TimeFieldFormat:  time.RFC3339,
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})
	r.Use(httplog.RequestLogger(logger, []string{"/ping"}))
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// API version 1.
	r.Route("/v1", func(r chi.Router) {
		r.Use(apiVersionCtx("v1"))
		r.Mount("/albums", albumRouter(albumHandler))
		r.Mount("/comments", commentRouter(commentHandler))
		r.Mount("/posts", postRouter(postHandler))
		r.Mount("/todos", todoRouter(todoHandler))
		r.Mount("/users", userRouter(userHandler))
	})
	http.ListenAndServe(":8000", r)
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "api.version", version))
			next.ServeHTTP(w, r)
		})
	}
}

func albumRouter(albumHandler *ah.AlbumHandler) http.Handler {
	r := chi.NewRouter()
	r.With(paginate).Get("/", albumHandler.GetAlbums)
	r.Post("/", albumHandler.CreateAlbum)
	r.Route("/{postID}", func(r chi.Router) {
		r.Get("/", albumHandler.GetAlbum)
		r.Put("/", albumHandler.UpdateAlbum)
		r.Delete("/", albumHandler.DeleteAlbum)
		r.Patch("/", albumHandler.PatchAlbum)
	})
	return r
}

func postRouter(postHandler *ph.PostHandler) http.Handler {
	r := chi.NewRouter()
	r.With(paginate).Get("/", postHandler.GetPosts)
	r.Post("/", postHandler.CreatePost)
	r.Route("/{postID}", func(r chi.Router) {
		r.Get("/", postHandler.GetPost)
		r.Put("/", postHandler.UpdatePost)
		r.Delete("/", postHandler.DeletePost)
		r.Patch("/", postHandler.PatchPost)
	})
	return r
}

func commentRouter(commentHandler *ch.CommentHandler) http.Handler {
	r := chi.NewRouter()
	r.With(paginate).Get("/", commentHandler.GetComments)
	r.Post("/", commentHandler.CreateComment)
	r.Route("/{postID}", func(r chi.Router) {
		r.Get("/", commentHandler.GetComment)
		r.Put("/", commentHandler.UpdateComment)
		r.Delete("/", commentHandler.DeleteComment)
		r.Patch("/", commentHandler.PatchComment)
	})
	return r
}

func todoRouter(todoHandler *th.TodoHandler) http.Handler {
	r := chi.NewRouter()
	r.With(paginate).Get("/", todoHandler.GetTodos)
	r.Post("/", todoHandler.CreateTodo)
	r.Route("/{postID}", func(r chi.Router) {
		r.Get("/", todoHandler.GetTodo)
		r.Put("/", todoHandler.UpdateTodo)
		r.Delete("/", todoHandler.DeleteTodo)
		r.Patch("/", todoHandler.PatchTodo)
	})
	return r
}
func userRouter(userHandler *uh.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.With(paginate).Get("/", userHandler.GetUsers)
	r.Post("/", userHandler.CreateUser)
	r.Route("/{postID}", func(r chi.Router) {
		r.Get("/", userHandler.GetUser)
		r.Put("/", userHandler.UpdateUser)
		r.Delete("/", userHandler.DeleteUser)
		r.Patch("/", userHandler.PatchUser)
	})
	return r
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}
