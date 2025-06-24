package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"muvi-discovery-app/internal/models"
)

const (
	TMDBBaseURL      = "https://api.themoviedb.org/3"
	TMDBImageBaseURL = "https://image.tmdb.org/t/p"
)

type TMDBService struct {
	apiKey     string
	httpClient *http.Client
}

func NewTMDBService(apiKey string) *TMDBService {
	return &TMDBService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *TMDBService) makeRequest(endpoint string, params url.Values) (*http.Response, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_key", s.apiKey)

	reqURL := fmt.Sprintf("%s%s?%s", TMDBBaseURL, endpoint, params.Encode())

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return resp, nil
}

// Movies
func (s *TMDBService) GetPopularMovies(page int) (*models.TMDBResponse[models.Movie], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))

	resp, err := s.makeRequest("/movie/popular", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.Movie]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetTopRatedMovies(page int) (*models.TMDBResponse[models.Movie], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))

	resp, err := s.makeRequest("/movie/top_rated", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.Movie]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetNowPlayingMovies(page int) (*models.TMDBResponse[models.Movie], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))

	resp, err := s.makeRequest("/movie/now_playing", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.Movie]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetMovieDetails(movieID int) (*models.MovieDetails, error) {
	endpoint := fmt.Sprintf("/movie/%d", movieID)

	resp, err := s.makeRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.MovieDetails
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetMovieCredits(movieID int) (*models.Credits, error) {
	endpoint := fmt.Sprintf("/movie/%d/credits", movieID)

	resp, err := s.makeRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.Credits
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetMovieVideos(movieID int) (*models.VideosResponse, error) {
	endpoint := fmt.Sprintf("/movie/%d/videos", movieID)

	resp, err := s.makeRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.VideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// TV Shows
func (s *TMDBService) GetPopularTVShows(page int) (*models.TMDBResponse[models.TVShow], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))

	resp, err := s.makeRequest("/tv/popular", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.TVShow]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetTopRatedTVShows(page int) (*models.TMDBResponse[models.TVShow], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))

	resp, err := s.makeRequest("/tv/top_rated", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.TVShow]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetTVShowDetails(tvID int) (*models.TVShowDetails, error) {
	endpoint := fmt.Sprintf("/tv/%d", tvID)
	params := url.Values{}
	params.Set("append_to_response", "external_ids")

	resp, err := s.makeRequest(endpoint, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TVShowDetails
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetTVShowVideos(tvID int) (*models.VideosResponse, error) {
	endpoint := fmt.Sprintf("/tv/%d/videos", tvID)

	resp, err := s.makeRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.VideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Search
func (s *TMDBService) SearchMovies(query string, page int) (*models.TMDBResponse[models.Movie], error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("page", strconv.Itoa(page))

	resp, err := s.makeRequest("/search/movie", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.Movie]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) SearchTVShows(query string, page int) (*models.TMDBResponse[models.TVShow], error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("page", strconv.Itoa(page))

	resp, err := s.makeRequest("/search/tv", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.TVShow]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Trending
func (s *TMDBService) GetTrendingMovies(timeWindow string) (*models.TMDBResponse[models.Movie], error) {
	endpoint := fmt.Sprintf("/trending/movie/%s", timeWindow)

	resp, err := s.makeRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.Movie]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetTrendingTVShows(timeWindow string) (*models.TMDBResponse[models.TVShow], error) {
	endpoint := fmt.Sprintf("/trending/tv/%s", timeWindow)

	resp, err := s.makeRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.TVShow]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Genres
func (s *TMDBService) GetMovieGenres() (*struct {
	Genres []models.Genre `json:"genres"`
}, error) {
	resp, err := s.makeRequest("/genre/movie/list", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Genres []models.Genre `json:"genres"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetTVGenres() (*struct {
	Genres []models.Genre `json:"genres"`
}, error) {
	resp, err := s.makeRequest("/genre/tv/list", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Genres []models.Genre `json:"genres"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Discover
func (s *TMDBService) DiscoverMovies(filters models.SearchFilters, page int) (*models.TMDBResponse[models.Movie], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))

	if filters.Genre != nil {
		params.Set("with_genres", strconv.Itoa(*filters.Genre))
	}
	if filters.Year != nil {
		params.Set("year", strconv.Itoa(*filters.Year))
	}
	if filters.Rating != nil {
		params.Set("vote_average.gte", fmt.Sprintf("%.1f", *filters.Rating))
	}
	if filters.SortBy != "" {
		sortDirection := "desc"
		if filters.SortOrder == "asc" {
			sortDirection = "asc"
		}
		params.Set("sort_by", fmt.Sprintf("%s.%s", filters.SortBy, sortDirection))
	}

	resp, err := s.makeRequest("/discover/movie", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.Movie]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) DiscoverTVShows(filters models.SearchFilters, page int) (*models.TMDBResponse[models.TVShow], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))

	if filters.Genre != nil {
		params.Set("with_genres", strconv.Itoa(*filters.Genre))
	}
	if filters.Year != nil {
		params.Set("first_air_date_year", strconv.Itoa(*filters.Year))
	}
	if filters.Rating != nil {
		params.Set("vote_average.gte", fmt.Sprintf("%.1f", *filters.Rating))
	}
	if filters.SortBy != "" {
		sortDirection := "desc"
		if filters.SortOrder == "asc" {
			sortDirection = "asc"
		}
		params.Set("sort_by", fmt.Sprintf("%s.%s", filters.SortBy, sortDirection))
	}

	resp, err := s.makeRequest("/discover/tv", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse[models.TVShow]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Utility functions
func (s *TMDBService) BuildImageURL(path *string, size string) string {
	if path == nil || *path == "" {
		return "/static/images/placeholder.jpg"
	}
	return fmt.Sprintf("%s/%s%s", TMDBImageBaseURL, size, *path)
}

func (s *TMDBService) GetPosterURL(path *string) string {
	return s.BuildImageURL(path, "w500")
}

func (s *TMDBService) GetBackdropURL(path *string) string {
	return s.BuildImageURL(path, "w1280")
}

func (s *TMDBService) GetThumbnailURL(path *string) string {
	return s.BuildImageURL(path, "w185")
}
