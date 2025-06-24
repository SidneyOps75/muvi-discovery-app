package main

import (
	"log"
	"net/http"
	"os"

	"muvi-discovery-app/internal/handlers"
	"muvi-discovery-app/internal/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Validate API keys
	tmdbKey := os.Getenv("TMDB_API_KEY")
	omdbKey := os.Getenv("OMDB_API_KEY")

	if tmdbKey == "" || tmdbKey == "demo_key" {
		log.Fatal("TMDB_API_KEY is required. Please set it in your .env file")
	}

	if omdbKey == "" || omdbKey == "demo_key" {
		log.Fatal("OMDB_API_KEY is required. Please set it in your .env file")
	}

	log.Printf("Loaded API keys - TMDB: %s..., OMDB: %s...", tmdbKey[:8], omdbKey[:8])

	// Initialize services
	tmdbService := services.NewTMDBService(tmdbKey)
	omdbService := services.NewOMDBService(omdbKey)

	// Initialize handlers
	h := handlers.NewHandler(tmdbService, omdbService)

	// Setup routes
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Routes
	r.HandleFunc("/", h.Home).Methods("GET")
	r.HandleFunc("/movies", h.Movies).Methods("GET")
	r.HandleFunc("/movies/{id}", h.MovieDetails).Methods("GET")
	r.HandleFunc("/tv", h.TVShows).Methods("GET")
	r.HandleFunc("/tv/{id}", h.TVShowDetails).Methods("GET")
	r.HandleFunc("/search", h.Search).Methods("GET")
	r.HandleFunc("/discover", h.Discover).Methods("GET")
	r.HandleFunc("/watchlist", h.Watchlist).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/search", h.APISearch).Methods("GET")
	api.HandleFunc("/watchlist", h.APIWatchlistAdd).Methods("POST")
	api.HandleFunc("/watchlist/{id}", h.APIWatchlistRemove).Methods("DELETE")
	api.HandleFunc("/watchlist/{id}/toggle", h.APIWatchlistToggle).Methods("PUT")
	api.HandleFunc("/movies/{id}/videos", h.APIMovieVideos).Methods("GET")
	api.HandleFunc("/tv/{id}/videos", h.APITVShowVideos).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
