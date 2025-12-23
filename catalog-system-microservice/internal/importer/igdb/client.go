package igdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Popularity  float64 `json:"popularity"`
}

type PopularMoviesResponse struct {
	Page       int     `json:"page"`
	Results    []Movie `json:"results"`
	TotalPages int     `json:"total_pages"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) getPopularMoviesPage(page int) (*PopularMoviesResponse, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/popular?api_key=%s&page=%d",
		c.apiKey,
		page,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb error: status %d", resp.StatusCode)
	}

	var data PopularMoviesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) GetFirst1000PopularMovies() ([]Movie, error) {
	const (
		moviesPerPage = 20
		targetCount   = 1000
		maxPages      = targetCount / moviesPerPage // 50
	)

	allMovies := make([]Movie, 0, targetCount)

	for page := 1; page <= maxPages; page++ {
		resp, err := c.getPopularMoviesPage(page)
		if err != nil {
			return nil, err
		}

		allMovies = append(allMovies, resp.Results...)

		if len(allMovies) >= targetCount {
			break
		}

		// Небольшая пауза, чтобы не словить rate-limit
		time.Sleep(250 * time.Millisecond)
	}

	return allMovies[:targetCount], nil
}

func (c *Client) GetPopularMoviesPage(page int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/popular?api_key=%s&page=%d",
		c.apiKey,
		page,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Results []map[string]interface{} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Results, nil
}
