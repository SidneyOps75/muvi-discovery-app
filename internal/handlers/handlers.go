package handlers

import (
	"encoding/json"
	"log"
	"muvi-discovery-app/internal/views"
	"net/http"
	"strconv"
	"strings"

	"muvi-discovery-app/internal/models"
	"muvi-discovery-app/internal/services"
	"muvi-discovery-app/web"

	"github.com/gorilla/mux"
)

type Handler struct {
	tmdbService      *services.TMDBService
	omdbService      *services.OMDBService
	watchlistService *services.WatchlistService
	templates        views.Template
}

func NewHandler(tmdbService *services.TMDBService, omdbService *services.OMDBService) *Handler {
	// Initialize watchlist service
	watchlistService := services.NewWatchlistService("data/watchlist.json")

	// Initialize templates
	tpl := views.Must(views.ParseFS(web.Templates, "templates/*.html"))

	return &Handler{
		tmdbService:      tmdbService,
		omdbService:      omdbService,
		watchlistService: watchlistService,
		templates:        tpl,
	}
}

type PageData struct {
	Title           string
	ContentTemplate string
	Movies          []models.Movie
	TVShows         []models.TVShow
	MovieDetails    *models.MovieDetails
	TVShowDetails   *models.TVShowDetails
	Credits         *models.Credits
	OMDBData        *models.OMDBMovie
	WatchlistItems  []models.WatchlistItem
	Genres          []models.Genre
	SearchQuery     string
	CurrentPage     int
	TotalPages      int
	Error           string
	IsInWatchlist   bool
	WatchlistCount  int
}

func (h *Handler) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data PageData) {
	// Add watchlist count to all pages
	data.WatchlistCount = h.watchlistService.GetItemCount()

	h.templates.Execute(w, r, name, data)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Muvi Discovery - Home"}

	log.Printf("Home handler called")

	// Get trending movies
	log.Printf("Fetching trending movies...")
	trendingMovies, err := h.tmdbService.GetTrendingMovies("week")
	if err != nil {
		log.Printf("Error fetching trending movies: %v", err)
		data.Error = "Failed to load trending movies"
	} else {
		log.Printf("Successfully fetched %d trending movies", len(trendingMovies.Results))
		// Limit to first 6 movies for homepage
		if len(trendingMovies.Results) > 6 {
			data.Movies = trendingMovies.Results[:6]
		} else {
			data.Movies = trendingMovies.Results
		}
	}

	// Get trending TV shows
	log.Printf("Fetching trending TV shows...")
	trendingTV, err := h.tmdbService.GetTrendingTVShows("week")
	if err != nil {
		log.Printf("Error fetching trending TV shows: %v", err)
		if data.Error == "" {
			data.Error = "Failed to load trending TV shows"
		}
	} else {
		log.Printf("Successfully fetched %d trending TV shows", len(trendingTV.Results))
		// Limit to first 6 shows for homepage
		if len(trendingTV.Results) > 6 {
			data.TVShows = trendingTV.Results[:6]
		} else {
			data.TVShows = trendingTV.Results
		}
	}

	log.Printf("Rendering template with %d movies and %d TV shows", len(data.Movies), len(data.TVShows))
	h.renderTemplate(w, r, "home.html", data)
}

func (h *Handler) Movies(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:           "Movies",
		ContentTemplate: "movies-content",
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "popular"
	}

	var moviesResp *models.TMDBResponse[models.Movie]
	var err error

	switch category {
	case "top_rated":
		moviesResp, err = h.tmdbService.GetTopRatedMovies(page)
	case "now_playing":
		moviesResp, err = h.tmdbService.GetNowPlayingMovies(page)
	default:
		moviesResp, err = h.tmdbService.GetPopularMovies(page)
	}

	if err != nil {
		log.Printf("Error fetching movies: %v", err)
		data.Error = "Failed to load movies"
	} else {
		data.Movies = moviesResp.Results
		data.CurrentPage = moviesResp.Page
		data.TotalPages = moviesResp.TotalPages
	}

	h.renderTemplate(w, r, "base.html", data)
}

func (h *Handler) MovieDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	data := PageData{
		Title:           "Movie Details",
		ContentTemplate: "movie-details-content",
	}

	// Get movie details
	movieDetails, err := h.tmdbService.GetMovieDetails(id)
	if err != nil {
		log.Printf("Error fetching movie details: %v", err)
		data.Error = "Failed to load movie details"
		h.renderTemplate(w, r, "base.html", data)
		return
	}

	data.MovieDetails = movieDetails
	data.Title = movieDetails.Title

	// Check if in watchlist
	data.IsInWatchlist = h.watchlistService.IsInWatchlist("movie", id)

	// Get credits
	credits, err := h.tmdbService.GetMovieCredits(id)
	if err != nil {
		log.Printf("Error fetching movie credits: %v", err)
	} else {
		data.Credits = credits
	}

	// Get OMDB data if IMDB ID is available
	if movieDetails.IMDBId != "" {
		omdbData, err := h.omdbService.GetMovieByIMDBID(movieDetails.IMDBId)
		if err != nil {
			log.Printf("Error fetching OMDB data: %v", err)
		} else {
			data.OMDBData = omdbData
		}
	}

	h.renderTemplate(w, r, "base.html", data)
}

func (h *Handler) TVShows(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:           "TV Shows",
		ContentTemplate: "tv-shows-content",
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "popular"
	}

	var tvResp *models.TMDBResponse[models.TVShow]
	var err error

	switch category {
	case "top_rated":
		tvResp, err = h.tmdbService.GetTopRatedTVShows(page)
	default:
		tvResp, err = h.tmdbService.GetPopularTVShows(page)
	}

	if err != nil {
		log.Printf("Error fetching TV shows: %v", err)
		data.Error = "Failed to load TV shows"
	} else {
		data.TVShows = tvResp.Results
		data.CurrentPage = tvResp.Page
		data.TotalPages = tvResp.TotalPages
	}

	h.renderTemplate(w, r, "base.html", data)
}

func (h *Handler) TVShowDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid TV show ID", http.StatusBadRequest)
		return
	}

	data := PageData{
		Title:           "TV Show Details",
		ContentTemplate: "tv-details-content",
	}

	// Get TV show details
	tvDetails, err := h.tmdbService.GetTVShowDetails(id)
	if err != nil {
		log.Printf("Error fetching TV show details: %v", err)
		data.Error = "Failed to load TV show details"
		h.renderTemplate(w, r, "base.html", data)
		return
	}

	data.TVShowDetails = tvDetails
	data.Title = tvDetails.Name

	// Check if in watchlist
	data.IsInWatchlist = h.watchlistService.IsInWatchlist("tv", id)

	// Get OMDB data if IMDB ID is available
	if tvDetails.ExternalIDs.IMDBID != "" {
		omdbData, err := h.omdbService.GetMovieByIMDBID(tvDetails.ExternalIDs.IMDBID)
		if err != nil {
			log.Printf("Error fetching OMDB data: %v", err)
		} else {
			data.OMDBData = omdbData
		}
	}

	h.renderTemplate(w, r, "base.html", data)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:           "Search",
		ContentTemplate: "search-content",
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		h.renderTemplate(w, r, "base.html", data)
		return
	}

	data.SearchQuery = query

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	mediaType := r.URL.Query().Get("type")
	if mediaType == "" {
		mediaType = "movie"
	}

	if mediaType == "tv" {
		tvResp, err := h.tmdbService.SearchTVShows(query, page)
		if err != nil {
			log.Printf("Error searching TV shows: %v", err)
			data.Error = "Failed to search TV shows"
		} else {
			data.TVShows = tvResp.Results
			data.CurrentPage = tvResp.Page
			data.TotalPages = tvResp.TotalPages
		}
	} else {
		moviesResp, err := h.tmdbService.SearchMovies(query, page)
		if err != nil {
			log.Printf("Error searching movies: %v", err)
			data.Error = "Failed to search movies"
		} else {
			data.Movies = moviesResp.Results
			data.CurrentPage = moviesResp.Page
			data.TotalPages = moviesResp.TotalPages
		}
	}

	h.renderTemplate(w, r, "base.html", data)
}

func (h *Handler) Discover(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:           "Discover",
		ContentTemplate: "discover-content",
	}

	// Get genres for filters
	movieGenres, err := h.tmdbService.GetMovieGenres()
	if err != nil {
		log.Printf("Error fetching genres: %v", err)
	} else {
		data.Genres = movieGenres.Genres
	}

	h.renderTemplate(w, r, "base.html", data)
}

func (h *Handler) Watchlist(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:           "My Watchlist",
		ContentTemplate: "watchlist-content",
	}

	filter := r.URL.Query().Get("filter")

	switch filter {
	case "watched":
		data.WatchlistItems = h.watchlistService.GetWatchedItems()
	case "unwatched":
		data.WatchlistItems = h.watchlistService.GetUnwatchedItems()
	default:
		data.WatchlistItems = h.watchlistService.GetAllItems()
	}

	h.renderTemplate(w, r, "base.html", data)
}

// API Handlers
func (h *Handler) APISearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "Query parameter required", http.StatusBadRequest)
		return
	}

	mediaType := r.URL.Query().Get("type")
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if mediaType == "tv" {
		tvResp, err := h.tmdbService.SearchTVShows(query, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tvResp)
	} else {
		moviesResp, err := h.tmdbService.SearchMovies(query, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(moviesResp)
	}
}

func (h *Handler) APIWatchlistAdd(w http.ResponseWriter, r *http.Request) {
	var item models.WatchlistItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.watchlistService.AddItem(item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) APIWatchlistRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	itemType := r.URL.Query().Get("type")
	if itemType == "" {
		http.Error(w, "Type parameter required", http.StatusBadRequest)
		return
	}

	if err := h.watchlistService.RemoveItem(itemType, id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) APIWatchlistToggle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	itemType := r.URL.Query().Get("type")
	if itemType == "" {
		http.Error(w, "Type parameter required", http.StatusBadRequest)
		return
	}

	if err := h.watchlistService.ToggleWatched(itemType, id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
