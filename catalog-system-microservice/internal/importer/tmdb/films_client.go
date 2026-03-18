package tmdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}

	// если задан прокси в env — используем его
	if proxyAddr := os.Getenv("HTTP_PROXY"); proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
		},
	}
}

func (c *Client) GetPopularMoviesPage(page int) (*PopularMoviesResponse, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/popular?api_key=%s&page=%d",
		c.apiKey, page,
	)
	result, err := doGet[PopularMoviesResponse](c, url)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetMovieDetails(id int) (map[string]interface{}, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%d?api_key=%s&append_to_response=credits,images,videos",
		id, c.apiKey,
	)
	return doGet[map[string]interface{}](c, url)
}

// GetFirst1000PopularMovies — с retry при rate limit
func (c *Client) GetFirst1000PopularMovies() ([]Movie, error) {
	const (
		moviesPerPage = 20
		targetCount   = 1000
		maxPages      = targetCount / moviesPerPage
	)

	result := make([]Movie, 0, targetCount)

	for page := 1; page <= maxPages; page++ {
		log.Printf("fetching page %d/%d", page, maxPages)

		resp, err := c.GetPopularMoviesPage(page)
		if err != nil {
			return nil, err
		}

		result = append(result, resp.Results...)
		time.Sleep(250 * time.Millisecond)
	}

	if len(result) > targetCount {
		result = result[:targetCount]
	}
	return result, nil
}

// doGet — универсальный GET с retry на 429
func doGet[T any](c *Client, url string) (T, error) {
	var zero T
	const maxRetries = 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := c.httpClient.Get(url)
		if err != nil {
			return zero, err
		}
		defer resp.Body.Close()

		// rate limit — ждём и повторяем
		if resp.StatusCode == http.StatusTooManyRequests {
			wait := time.Duration(attempt*attempt) * time.Second
			log.Printf("rate limited, waiting %s (attempt %d/%d)", wait, attempt, maxRetries)
			time.Sleep(wait)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return zero, fmt.Errorf("tmdb error: status %d", resp.StatusCode)
		}

		var data T
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return zero, err
		}
		return data, nil
	}

	return zero, fmt.Errorf("max retries exceeded for %s", url)
}
