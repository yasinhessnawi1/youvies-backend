package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers(r *mux.Router) {
	// Movies Endpoints
	r.HandleFunc("/youvies/v1/movies", AuthMiddleware(http.HandlerFunc(GetMovies), "user").ServeHTTP).Methods("GET", "OPTIONS")
	r.HandleFunc("/youvies/v1/movies", AuthMiddleware(http.HandlerFunc(CreateMovie), "admin").ServeHTTP).Methods("POST", "OPTIONS")
	r.HandleFunc("/youvies/v1/movies/{id}", AuthMiddleware(http.HandlerFunc(UpdateMovie), "admin").ServeHTTP).Methods("PUT", "OPTIONS")
	r.HandleFunc("/youvies/v1/movies/{id}", AuthMiddleware(http.HandlerFunc(DeleteMovie), "admin").ServeHTTP).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/youvies/v1/movies/search", AuthMiddleware(http.HandlerFunc(SearchMovies), "user").ServeHTTP).Methods("GET", "OPTIONS")

	// Anime Shows Endpoints
	r.HandleFunc("/youvies/v1/animeshows", AuthMiddleware(http.HandlerFunc(GetAnimeShows), "user").ServeHTTP).Methods("GET", "OPTIONS")
	r.HandleFunc("/youvies/v1/animeshows", AuthMiddleware(http.HandlerFunc(CreateAnimeShow), "admin").ServeHTTP).Methods("POST", "OPTIONS")
	r.HandleFunc("/youvies/v1/animeshows/{id}", AuthMiddleware(http.HandlerFunc(UpdateAnimeShow), "admin").ServeHTTP).Methods("PUT", "OPTIONS")
	r.HandleFunc("/youvies/v1/animeshows/{id}", AuthMiddleware(http.HandlerFunc(DeleteAnimeShow), "admin").ServeHTTP).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/youvies/v1/animeshows/search", AuthMiddleware(http.HandlerFunc(SearchAnimeShows), "user").ServeHTTP).Methods("GET", "OPTIONS")

	// Anime Movies Endpoints
	r.HandleFunc("/youvies/v1/animemovies", AuthMiddleware(http.HandlerFunc(GetAnimeMovies), "user").ServeHTTP).Methods("GET", "OPTIONS")
	r.HandleFunc("/youvies/v1/animemovies", AuthMiddleware(http.HandlerFunc(CreateAnimeMovie), "admin").ServeHTTP).Methods("POST", "OPTIONS")
	r.HandleFunc("/youvies/v1/animemovies/{id}", AuthMiddleware(http.HandlerFunc(UpdateAnimeMovie), "admin").ServeHTTP).Methods("PUT", "OPTIONS")
	r.HandleFunc("/youvies/v1/animemovies/{id}", AuthMiddleware(http.HandlerFunc(DeleteAnimeMovie), "admin").ServeHTTP).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/youvies/v1/animemovies/search", AuthMiddleware(http.HandlerFunc(SearchAnimeMovies), "user").ServeHTTP).Methods("GET", "OPTIONS")

	// Shows Endpoints
	r.HandleFunc("/youvies/v1/shows", AuthMiddleware(http.HandlerFunc(GetShows), "user").ServeHTTP).Methods("GET", "OPTIONS")
	r.HandleFunc("/youvies/v1/shows", AuthMiddleware(http.HandlerFunc(CreateShow), "admin").ServeHTTP).Methods("POST", "OPTIONS")
	r.HandleFunc("/youvies/v1/shows/{id}", AuthMiddleware(http.HandlerFunc(UpdateShow), "admin").ServeHTTP).Methods("PUT", "OPTIONS")
	r.HandleFunc("/youvies/v1/shows/{id}", AuthMiddleware(http.HandlerFunc(DeleteShow), "admin").ServeHTTP).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/youvies/v1/shows/search", AuthMiddleware(http.HandlerFunc(SearchShows), "user").ServeHTTP).Methods("GET", "OPTIONS")

	// User Endpoints
	r.HandleFunc("/youvies/v1/api/register", RegisterUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/youvies/v1/api/login", LoginUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/youvies/v1/api/logout", AuthMiddleware(http.HandlerFunc(LogoutUser), "user").ServeHTTP).Methods("POST", "OPTIONS")
	r.HandleFunc("/youvies/v1/api/user", AuthMiddleware(http.HandlerFunc(EditUser), "user").ServeHTTP).Methods("PUT", "OPTIONS")
}
