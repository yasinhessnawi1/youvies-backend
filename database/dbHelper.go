package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"sort"
	"youvies-backend/models"
)

func marshalJSON(data interface{}) ([]byte, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %v", err)
	}
	return result, nil
}

func unmarshalJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("error unmarshaling data: %v", err)
	}
	return nil
}

func insertAttributes(attributes models.Attributes) (int64, error) {
	titles, err := marshalJSON(attributes.Titles)
	if err != nil {
		return 0, err
	}

	posterImage, err := marshalJSON(attributes.PosterImage)
	if err != nil {
		return 0, err
	}

	coverImage, err := marshalJSON(attributes.CoverImage)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO attributes (
        created_at, updated_at, slug, synopsis, description, titles, canonical_title, average_rating, favorites_count,
        start_date, end_date, next_release, popularity_rank, rating_rank, age_rating, age_rating_guide, subtype, status,
        poster_image, cover_image, episode_count, episode_length, youtube_video_id, show_type
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
    ) RETURNING id`

	var id int64
	err = DB.QueryRow(query,
		attributes.CreatedAt, attributes.UpdatedAt, attributes.Slug, attributes.Synopsis, attributes.Description, titles, attributes.CanonicalTitle,
		attributes.AverageRating, attributes.FavoritesCount, attributes.StartDate, attributes.EndDate, attributes.NextRelease, attributes.PopularityRank,
		attributes.RatingRank, attributes.AgeRating, attributes.AgeRatingGuide, attributes.Subtype, attributes.Status, posterImage, coverImage,
		attributes.EpisodeCount, attributes.EpisodeLength, attributes.YoutubeVideoId, attributes.ShowType).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting attributes: %v", err)
	}

	return id, nil
}

func insertRelationships(relationships models.Relationships) (int64, error) {
	query := `INSERT INTO relationships (self, related) VALUES ($1, $2) RETURNING id`

	var id int64
	err := DB.QueryRow(query, relationships.Genres.Links.Self, relationships.Genres.Links.Related).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting relationships: %v", err)
	}
	return id, nil
}

func InsertEpisode(episode models.Episode) (int64, error) {
	titles, err := marshalJSON(episode.Titles)
	if err != nil {
		return 0, err
	}

	thumbnail, err := marshalJSON(episode.Thumbnail)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO episodes (
        id, anime_show_id, created_at, updated_at, synopsis, description, titles, canonical_title, season_number, number,
        relative_number, airdate, length, thumbnail
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
    ) RETURNING id`

	var id int64
	err = DB.QueryRow(query, episode.ID,
		episode.AnimeShowID, episode.CreatedAt, episode.UpdatedAt, episode.Synopsis, episode.Description, titles, episode.CanonicalTitle,
		episode.SeasonNumber, episode.Number, episode.RelativeNumber, episode.Airdate, episode.Length, thumbnail).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting episode: %v", err)
	}

	return id, nil
}

func insertAnimeShow(animeShow models.AnimeShow) (string, error) {
	attributesID, err := insertAttributes(animeShow.Attributes)
	if err != nil {
		return "", err
	}
	relationshipsID, err := insertRelationships(animeShow.Relationships)
	if err != nil {
		return "", err
	}
	genres, err := marshalJSON(animeShow.Genres)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO anime_shows (
        id, title, attributes_id, relationships_id, genres
    ) VALUES (
        $1, $2, $3, $4, $5
    ) RETURNING id`

	var id string
	err = DB.QueryRow(query, animeShow.ID, animeShow.Title, attributesID, relationshipsID, genres).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("error inserting anime show: %v", err)
	}

	return id, nil
}

func insertAnimeMovie(animeMovie models.AnimeMovie) (string, error) {
	attributesID, err := insertAttributes(animeMovie.Attributes)
	if err != nil {
		return "", err
	}
	relationshipsID, err := insertRelationships(animeMovie.Relationships)
	if err != nil {
		return "", err
	}
	genres, err := marshalJSON(animeMovie.Genres)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO anime_movies (
		id, title, attributes_id, relationships_id, genres
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING id`

	var id string
	err = DB.QueryRow(query, animeMovie.ID, animeMovie.Title, attributesID, relationshipsID, genres).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("error inserting anime movie: %v", err)
	}

	return id, nil
}

func insertMovie(movie models.Movie) error {
	genres, err := marshalJSON(movie.Genres)
	if err != nil {
		return fmt.Errorf("error marshaling genres: %v", err)
	}

	query := `INSERT INTO movie (
		id, original_language, original_title, overview, popularity, poster_path, release_date, title, vote_average, vote_count, backdrop_path, genres, last_updated
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)`

	_, err = DB.Exec(query, movie.ID, movie.OriginalLanguage, movie.OriginalTitle, movie.Overview, movie.Popularity, movie.PosterPath, movie.ReleaseDate, movie.Title, movie.VoteAverage, movie.VoteCount, movie.BackdropPath, genres, movie.LastUpdated)
	if err != nil {
		return fmt.Errorf("error inserting movie: %v", err)
	}
	return nil
}

func insertShow(show models.Show) error {
	genres, err := marshalJSON(show.Genres)
	if err != nil {
		return fmt.Errorf("error marshaling genres: %v", err)
	}

	seasonsInfo, err := marshalJSON(show.SeasonsInfo)
	if err != nil {
		return fmt.Errorf("error marshaling seasons info: %v", err)
	}

	query := `INSERT INTO show (
		id, title, overview, poster_path, first_air_date, genres, vote_average, vote_count, original_language, popularity, backdrop_path, last_updated, seasons_info
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)`

	_, err = DB.Exec(query, show.ID, show.Title, show.Overview, show.PosterPath, show.FirstAirDate, genres, show.VoteAverage, show.VoteCount, show.OriginalLanguage, show.Popularity, show.BackdropPath, show.LastUpdated, seasonsInfo)
	if err != nil {
		return fmt.Errorf("error inserting show: %v", err)
	}
	return nil
}

func insertUser(user models.User) error {
	query := `INSERT INTO users (
		id, username, password, email, role, active, created, updated, avatar, favorites, friends, rooms, watched
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)`

	_, err := DB.Exec(query, user.ID, user.Username, user.Password, user.Email, user.Role, user.Active, user.Created, user.Updated, user.Avatar, pq.Array(user.Favorites), pq.Array(user.Friends), pq.Array(user.Rooms), pq.Array(user.Watched))
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func getAttributes(attributesID int64) (models.Attributes, error) {
	var attributes models.Attributes
	var titles, posterImage, coverImage []byte
	query := `SELECT id, created_at, updated_at, slug, synopsis, description, titles, canonical_title, average_rating, favorites_count, 
              start_date, end_date, next_release, popularity_rank, rating_rank, age_rating, age_rating_guide, subtype, status, poster_image, cover_image, 
              episode_count, episode_length, youtube_video_id, show_type 
              FROM attributes WHERE id=$1`
	err := DB.QueryRow(query, attributesID).Scan(
		&attributes.ID, &attributes.CreatedAt, &attributes.UpdatedAt, &attributes.Slug, &attributes.Synopsis, &attributes.Description, &titles,
		&attributes.CanonicalTitle, &attributes.AverageRating, &attributes.FavoritesCount, &attributes.StartDate, &attributes.EndDate, &attributes.NextRelease,
		&attributes.PopularityRank, &attributes.RatingRank, &attributes.AgeRating, &attributes.AgeRatingGuide, &attributes.Subtype, &attributes.Status, &posterImage,
		&coverImage, &attributes.EpisodeCount, &attributes.EpisodeLength, &attributes.YoutubeVideoId, &attributes.ShowType,
	)
	if err != nil {
		return attributes, fmt.Errorf("error getting attributes: %v", err)
	}

	unmarshalJSON(titles, &attributes.Titles)
	unmarshalJSON(posterImage, &attributes.PosterImage)
	unmarshalJSON(coverImage, &attributes.CoverImage)

	return attributes, nil
}

func getRelationships(relationshipsID int64) (models.Relationships, error) {
	var relationships models.Relationships
	query := `SELECT id, self, related FROM relationships WHERE id=$1`
	err := DB.QueryRow(query, relationshipsID).Scan(
		&relationships.ID, &relationships.Genres.Links.Self, &relationships.Genres.Links.Related,
	)
	if err != nil {
		return relationships, fmt.Errorf("error getting relationships: %v", err)
	}
	return relationships, nil
}

func getEpisode(episodeID int64) (models.Episode, error) {
	var episode models.Episode
	var titles, thumbnail []byte
	query := `SELECT id, anime_show_id, created_at, updated_at, synopsis, description, titles, canonical_title, season_number, number, 
              relative_number, airdate, length, thumbnail 
              FROM episodes WHERE id=$1`
	err := DB.QueryRow(query, episodeID).Scan(
		&episode.ID, &episode.AnimeShowID, &episode.CreatedAt, &episode.UpdatedAt, &episode.Synopsis, &episode.Description, &titles,
		&episode.CanonicalTitle, &episode.SeasonNumber, &episode.Number, &episode.RelativeNumber, &episode.Airdate, &episode.Length, &thumbnail,
	)
	if err != nil {
		return episode, fmt.Errorf("error getting episode: %v", err)
	}

	unmarshalJSON(titles, &episode.Titles)
	unmarshalJSON(thumbnail, &episode.Thumbnail)

	return episode, nil
}

func getAnimeShow(animeShowID string) (models.AnimeShow, error) {
	var animeShow models.AnimeShow
	var genres []byte
	var attributesID, relationshipsID int64
	query := `SELECT id, title, attributes_id, relationships_id, genres FROM anime_shows WHERE id=$1`
	err := DB.QueryRow(query, animeShowID).Scan(
		&animeShow.ID, &animeShow.Title, &attributesID, &relationshipsID, &genres,
	)
	if err != nil {
		return animeShow, fmt.Errorf("error getting anime show: %v", err)
	}

	animeShow.Attributes, err = getAttributes(attributesID)
	if err != nil {
		return animeShow, err
	}

	animeShow.Relationships, err = getRelationships(relationshipsID)
	if err != nil {
		return animeShow, err
	}

	unmarshalJSON(genres, &animeShow.Genres)

	return animeShow, nil
}

func getAnimeMovie(animeMovieID string) (models.AnimeMovie, error) {
	var animeMovie models.AnimeMovie
	var genres []byte
	var attributesID, relationshipsID int64
	query := `SELECT id, title, attributes_id, relationships_id, genres FROM anime_movies WHERE id=$1`
	err := DB.QueryRow(query, animeMovieID).Scan(
		&animeMovie.ID, &animeMovie.Title, &attributesID, &relationshipsID, &genres,
	)
	if err != nil {
		return animeMovie, fmt.Errorf("error getting anime movie: %v", err)
	}

	animeMovie.Attributes, err = getAttributes(attributesID)
	if err != nil {
		return animeMovie, err
	}

	animeMovie.Relationships, err = getRelationships(relationshipsID)
	if err != nil {
		return animeMovie, err
	}

	unmarshalJSON(genres, &animeMovie.Genres)

	return animeMovie, nil
}

func getMovie(movieID string) (models.Movie, error) {
	var movie models.Movie
	var genres []byte
	query := `SELECT id, original_language, original_title, overview, popularity, poster_path, release_date, title, vote_average, vote_count, backdrop_path, genres, last_updated 
              FROM movie WHERE id=$1`
	err := DB.QueryRow(query, movieID).Scan(
		&movie.ID, &movie.OriginalLanguage, &movie.OriginalTitle, &movie.Overview, &movie.Popularity, &movie.PosterPath, &movie.ReleaseDate, &movie.Title,
		&movie.VoteAverage, &movie.VoteCount, &movie.BackdropPath, &genres, &movie.LastUpdated,
	)
	if err != nil {
		return movie, fmt.Errorf("error getting movie: %v", err)
	}

	unmarshalJSON(genres, &movie.Genres)

	return movie, nil
}

func getShow(showID string) (models.Show, error) {
	var show models.Show
	var genres, seasonsInfo []byte
	query := `SELECT id, title, overview, poster_path, first_air_date, genres, vote_average, vote_count, original_language, popularity, backdrop_path, last_updated, seasons_info 
              FROM show WHERE id=$1`
	err := DB.QueryRow(query, showID).Scan(
		&show.ID, &show.Title, &show.Overview, &show.PosterPath, &show.FirstAirDate, &genres, &show.VoteAverage, &show.VoteCount, &show.OriginalLanguage,
		&show.Popularity, &show.BackdropPath, &show.LastUpdated, &seasonsInfo,
	)
	if err != nil {
		return show, fmt.Errorf("error getting show: %v", err)
	}

	unmarshalJSON(genres, &show.Genres)
	unmarshalJSON(seasonsInfo, &show.SeasonsInfo)

	return show, nil
}

func getUser(userID string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password, email, role, active, created, updated, avatar, favorites, friends, rooms, watched 
              FROM users WHERE id=$1`
	err := DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.Active, &user.Created, &user.Updated, &user.Avatar,
		pq.Array(&user.Favorites), pq.Array(&user.Friends), pq.Array(&user.Rooms), pq.Array(&user.Watched),
	)
	if err != nil {
		return user, fmt.Errorf("error getting user: %v", err)
	}

	return user, nil
}

func getEpisodesByAnimeShowID(animeShowID string) ([]models.Episode, error) {
	var episodes []models.Episode
	rows, err := DB.Query(`SELECT id, anime_show_id, created_at, updated_at, synopsis, description, titles, canonical_title, season_number, number, 
                           relative_number, airdate, length, thumbnail FROM episodes WHERE anime_show_id=$1`, animeShowID)
	if err != nil {
		return episodes, fmt.Errorf("error getting episodes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var episode models.Episode
		var titles, thumbnail []byte
		err := rows.Scan(&episode.ID, &episode.AnimeShowID, &episode.CreatedAt, &episode.UpdatedAt, &episode.Synopsis, &episode.Description, &titles,
			&episode.CanonicalTitle, &episode.SeasonNumber, &episode.Number, &episode.RelativeNumber, &episode.Airdate, &episode.Length, &thumbnail)
		if err != nil {
			return episodes, fmt.Errorf("error scanning episode row: %v", err)
		}

		unmarshalJSON(titles, &episode.Titles)
		unmarshalJSON(thumbnail, &episode.Thumbnail)
		episodes = append(episodes, episode)
	}

	return episodes, nil
}

func scanRows(rows *sql.Rows, results interface{}) error {
	switch v := results.(type) {
	case *[]models.Movie:
		var movies []models.Movie
		for rows.Next() {
			var movie models.Movie
			var genres []byte
			err := rows.Scan(&movie.ID, &movie.Overview, &movie.PosterPath, &movie.ReleaseDate, &movie.Title, &movie.VoteAverage, &movie.BackdropPath, &movie.Popularity, &genres)
			if err != nil {
				return fmt.Errorf("error scanning row: %v", err)
			}
			unmarshalJSON(genres, &movie.Genres)
			movies = append(movies, movie)
		}
		// Sort movies by ReleaseDate in descending order
		sort.SliceStable(movies, func(i, j int) bool {
			return movies[i].ReleaseDate > movies[j].ReleaseDate
		})
		*v = movies

	case *[]models.Show:
		var shows []models.Show
		for rows.Next() {
			var show models.Show
			var genres, seasonsInfo []byte
			err := rows.Scan(&show.ID, &show.Title, &show.Overview, &show.PosterPath, &show.FirstAirDate, &genres, &show.VoteAverage, &show.VoteCount, &show.OriginalLanguage, &show.Popularity, &show.BackdropPath, &show.LastUpdated, &seasonsInfo)
			if err != nil {
				return fmt.Errorf("error scanning row: %v", err)
			}
			unmarshalJSON(genres, &show.Genres)
			unmarshalJSON(seasonsInfo, &show.SeasonsInfo)
			shows = append(shows, show)
		}
		// Sort shows by FirstAirDate in descending order
		sort.SliceStable(shows, func(i, j int) bool {
			return shows[i].FirstAirDate > shows[j].FirstAirDate
		})
		*v = shows

	case *[]models.AnimeTiny:
		var animeShows []models.AnimeTiny
		for rows.Next() {
			var animeShow models.AnimeTiny
			var genres, posterPath, coverImage []byte
			err := rows.Scan(&animeShow.ID, &animeShow.Title, &genres, &animeShow.AttributeID,
				&animeShow.Slug, &animeShow.Synopsis,
				&animeShow.Description,
				&animeShow.CanonicalTitle,
				&animeShow.AverageRating,
				&animeShow.PopularityRank,
				&animeShow.RatingRank,
				&animeShow.SubType,
				&animeShow.EpisodeCount,
				&posterPath,
				&coverImage)
			if err != nil {
				return fmt.Errorf("error scanning row: %v", err)
			}
			unmarshalJSON(genres, &animeShow.Genres)
			unmarshalJSON(posterPath, &animeShow.PosterImage)
			unmarshalJSON(coverImage, &animeShow.CoverImage)
			animeShows = append(animeShows, animeShow)
		}
		*v = animeShows

	case *[]models.User:
		var users []models.User
		for rows.Next() {
			var user models.User
			err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.Active, &user.Created, &user.Updated, &user.Avatar, pq.Array(&user.Favorites), pq.Array(&user.Friends), pq.Array(&user.Rooms), pq.Array(&user.Watched))
			if err != nil {
				return fmt.Errorf("error scanning row: %v", err)
			}
			users = append(users, user)
		}
		*v = users

	default:
		return fmt.Errorf("unsupported result type")
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %v", err)
	}

	return nil
}
