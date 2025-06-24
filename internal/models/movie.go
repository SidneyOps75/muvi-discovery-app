package models

import "time"

// Movie represents a movie from TMDB API
type Movie struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	Overview         string  `json:"overview"`
	PosterPath       *string `json:"poster_path"`
	BackdropPath     *string `json:"backdrop_path"`
	ReleaseDate      string  `json:"release_date"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	GenreIDs         []int   `json:"genre_ids"`
	Adult            bool    `json:"adult"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Popularity       float64 `json:"popularity"`
	Video            bool    `json:"video"`
}

// MovieDetails represents detailed movie information
type MovieDetails struct {
	Movie
	BelongsToCollection interface{}         `json:"belongs_to_collection"`
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	IMDBId              string              `json:"imdb_id"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
}

// TVShow represents a TV show from TMDB API
type TVShow struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	Overview         string   `json:"overview"`
	PosterPath       *string  `json:"poster_path"`
	BackdropPath     *string  `json:"backdrop_path"`
	FirstAirDate     string   `json:"first_air_date"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	GenreIDs         []int    `json:"genre_ids"`
	Adult            bool     `json:"adult"`
	OriginalLanguage string   `json:"original_language"`
	OriginalName     string   `json:"original_name"`
	Popularity       float64  `json:"popularity"`
	OriginCountry    []string `json:"origin_country"`
}

// TVShowDetails represents detailed TV show information
type TVShowDetails struct {
	TVShow
	CreatedBy           []interface{}       `json:"created_by"`
	EpisodeRunTime      []int               `json:"episode_run_time"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	InProduction        bool                `json:"in_production"`
	Languages           []string            `json:"languages"`
	LastAirDate         string              `json:"last_air_date"`
	LastEpisodeToAir    interface{}         `json:"last_episode_to_air"`
	NextEpisodeToAir    interface{}         `json:"next_episode_to_air"`
	Networks            []interface{}       `json:"networks"`
	NumberOfEpisodes    int                 `json:"number_of_episodes"`
	NumberOfSeasons     int                 `json:"number_of_seasons"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	Seasons             []interface{}       `json:"seasons"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Type                string              `json:"type"`
	ExternalIDs         ExternalIDs         `json:"external_ids"`
}

// ExternalIDs represents a list of external IDs (e.g., IMDB, TVDB)
type ExternalIDs struct {
	IMDBID      string `json:"imdb_id"`
	TVDBID      int    `json:"tvdb_id"`
	FacebookID  string `json:"facebook_id"`
	InstagramID string `json:"instagram_id"`
	TwitterID   string `json:"twitter_id"`
}

// Genre represents a movie/TV genre
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ProductionCompany represents a production company
type ProductionCompany struct {
	ID            int     `json:"id"`
	LogoPath      *string `json:"logo_path"`
	Name          string  `json:"name"`
	OriginCountry string  `json:"origin_country"`
}

// ProductionCountry represents a production country
type ProductionCountry struct {
	ISO31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

// SpokenLanguage represents a spoken language
type SpokenLanguage struct {
	EnglishName string `json:"english_name"`
	ISO6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
}

// TMDBResponse represents a paginated response from TMDB
type TMDBResponse[T any] struct {
	Page         int `json:"page"`
	Results      []T `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// OMDBMovie represents a movie from OMDB API
type OMDBMovie struct {
	Title      string   `json:"Title"`
	Year       string   `json:"Year"`
	Rated      string   `json:"Rated"`
	Released   string   `json:"Released"`
	Runtime    string   `json:"Runtime"`
	Genre      string   `json:"Genre"`
	Director   string   `json:"Director"`
	Writer     string   `json:"Writer"`
	Actors     string   `json:"Actors"`
	Plot       string   `json:"Plot"`
	Language   string   `json:"Language"`
	Country    string   `json:"Country"`
	Awards     string   `json:"Awards"`
	Poster     string   `json:"Poster"`
	Ratings    []Rating `json:"Ratings"`
	Metascore  string   `json:"Metascore"`
	IMDBRating string   `json:"imdbRating"`
	IMDBVotes  string   `json:"imdbVotes"`
	IMDBID     string   `json:"imdbID"`
	Type       string   `json:"Type"`
	DVD        string   `json:"DVD"`
	BoxOffice  string   `json:"BoxOffice"`
	Production string   `json:"Production"`
	Website    string   `json:"Website"`
	Response   string   `json:"Response"`
}

// Rating represents a rating from various sources
type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

// WatchlistItem represents an item in the user's watchlist
type WatchlistItem struct {
	ID          int        `json:"id"`
	Type        string     `json:"type"` // "movie" or "tv"
	Title       string     `json:"title"`
	PosterPath  *string    `json:"poster_path"`
	ReleaseDate string     `json:"release_date"`
	VoteAverage float64    `json:"vote_average"`
	Watched     bool       `json:"watched"`
	AddedAt     time.Time  `json:"added_at"`
	WatchedAt   *time.Time `json:"watched_at,omitempty"`
}

// SearchFilters represents search and discovery filters
type SearchFilters struct {
	Genre     *int     `json:"genre,omitempty"`
	Year      *int     `json:"year,omitempty"`
	Rating    *float64 `json:"rating,omitempty"`
	SortBy    string   `json:"sort_by,omitempty"`
	SortOrder string   `json:"sort_order,omitempty"`
}

// Credits represents cast and crew information
type Credits struct {
	ID   int          `json:"id"`
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

// CastMember represents a cast member
type CastMember struct {
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	ID                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
	CastID             int     `json:"cast_id"`
	Character          string  `json:"character"`
	CreditID           string  `json:"credit_id"`
	Order              int     `json:"order"`
}

// CrewMember represents a crew member
type CrewMember struct {
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	ID                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        *string `json:"profile_path"`
	CreditID           string  `json:"credit_id"`
	Department         string  `json:"department"`
	Job                string  `json:"job"`
}
