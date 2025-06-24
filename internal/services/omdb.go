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

const OMDBBaseURL = "https://www.omdbapi.com"

type OMDBService struct {
	apiKey     string
	httpClient *http.Client
}

func NewOMDBService(apiKey string) *OMDBService {
	return &OMDBService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *OMDBService) makeRequest(params url.Values) (*http.Response, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apikey", s.apiKey)

	reqURL := fmt.Sprintf("%s?%s", OMDBBaseURL, params.Encode())
	
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

func (s *OMDBService) GetMovieByIMDBID(imdbID string) (*models.OMDBMovie, error) {
	params := url.Values{}
	params.Set("i", imdbID)
	params.Set("plot", "full")

	resp, err := s.makeRequest(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.OMDBMovie
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Response == "False" {
		return nil, fmt.Errorf("movie not found")
	}

	return &result, nil
}

func (s *OMDBService) GetMovieByTitle(title string, year *int) (*models.OMDBMovie, error) {
	params := url.Values{}
	params.Set("t", title)
	params.Set("plot", "full")
	
	if year != nil {
		params.Set("y", strconv.Itoa(*year))
	}

	resp, err := s.makeRequest(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.OMDBMovie
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Response == "False" {
		return nil, fmt.Errorf("movie not found")
	}

	return &result, nil
}

func (s *OMDBService) SearchMovies(title string, page int) (*struct {
	Search       []models.OMDBMovie `json:"Search"`
	TotalResults string             `json:"totalResults"`
	Response     string             `json:"Response"`
}, error) {
	params := url.Values{}
	params.Set("s", title)
	params.Set("page", strconv.Itoa(page))
	params.Set("type", "movie")

	resp, err := s.makeRequest(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Search       []models.OMDBMovie `json:"Search"`
		TotalResults string             `json:"totalResults"`
		Response     string             `json:"Response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Response == "False" {
		return nil, fmt.Errorf("no results found")
	}

	return &result, nil
}

func (s *OMDBService) GetRatings(imdbID string) (*struct {
	IMDBRating string          `json:"imdbRating"`
	IMDBVotes  string          `json:"imdbVotes"`
	Ratings    []models.Rating `json:"Ratings"`
}, error) {
	movie, err := s.GetMovieByIMDBID(imdbID)
	if err != nil {
		return nil, err
	}

	return &struct {
		IMDBRating string          `json:"imdbRating"`
		IMDBVotes  string          `json:"imdbVotes"`
		Ratings    []models.Rating `json:"Ratings"`
	}{
		IMDBRating: movie.IMDBRating,
		IMDBVotes:  movie.IMDBVotes,
		Ratings:    movie.Ratings,
	}, nil
}

// Helper methods to extract specific ratings
func (s *OMDBService) GetRatingBySource(ratings []models.Rating, source string) *string {
	for _, rating := range ratings {
		if rating.Source == source {
			return &rating.Value
		}
	}
	return nil
}

func (s *OMDBService) GetIMDBRating(ratings []models.Rating) *string {
	return s.GetRatingBySource(ratings, "Internet Movie Database")
}

func (s *OMDBService) GetRottenTomatoesRating(ratings []models.Rating) *string {
	return s.GetRatingBySource(ratings, "Rotten Tomatoes")
}

func (s *OMDBService) GetMetacriticRating(ratings []models.Rating) *string {
	return s.GetRatingBySource(ratings, "Metacritic")
}