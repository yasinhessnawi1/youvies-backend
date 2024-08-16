package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Movies Endpoints
	/*
		r.GET("/youvies/v1/movies", AuthMiddleware("user"), GetMovies)
		r.GET("/youvies/v1/movies/search", AuthMiddleware("user"), SearchMovies)
		r.GET("/youvies/v1/movies/genre/:genre", AuthMiddleware("user"), GetMoviesByGenre)
		r.GET("/youvies/v1/movies/:id", AuthMiddleware("user"), GetMovieByID)

		// Anime Shows Endpoints
		r.GET("/youvies/v1/animeshows", AuthMiddleware("user"), GetAnimeShows)
		r.GET("/youvies/v1/animeshows/search", AuthMiddleware("user"), SearchAnimeShows)
		r.GET("/youvies/v1/animeshows/genre/:genre", AuthMiddleware("user"), GetAnimeShowsByGenre)
		r.GET("/youvies/v1/animeshows/:id", AuthMiddleware("user"), GetAnimeShowByID)

		// Anime Movies Endpoints
		r.GET("/youvies/v1/animemovies", AuthMiddleware("user"), GetAnimeMovies)
		r.GET("/youvies/v1/animemovies/search", AuthMiddleware("user"), SearchAnimeMovies)
		r.GET("/youvies/v1/animemovies/genre/:genre", AuthMiddleware("user"), GetAnimeMoviesByGenre)
		r.GET("/youvies/v1/animemovies/:id", AuthMiddleware("user"), GetAnimeMovieByID)

		// Shows Endpoints
		r.GET("/youvies/v1/shows", AuthMiddleware("user"), GetShows)
		r.GET("/youvies/v1/shows/search", AuthMiddleware("user"), SearchShows)
		r.GET("/youvies/v1/shows/genre/:genre", AuthMiddleware("user"), GetShowsByGenre)
		r.GET("/youvies/v1/shows/:id", AuthMiddleware("user"), GetShowByID)


	*/
	// User Endpoints
	r.POST("/youvies/v1/api/register", RegisterUser)
	r.POST("/youvies/v1/api/login", LoginUser)
	r.POST("/youvies/v1/api/logout", AuthMiddleware("user"), LogoutUser)
	r.PUT("/youvies/v1/api/user", AuthMiddleware("user"), EditUser)

}
