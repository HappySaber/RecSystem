package tmdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

/*
	=====================
	DTO для TMDB
	=====================
*/

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

/*
	=====================
	Клиент
	=====================
*/

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	transport := &http.Transport{
		Proxy: nil,
	}

	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		},
	}
}

/*
	=====================
	Популярные фильмы (1 страница)
	=====================
*/

func (c *Client) GetPopularMoviesPage(page int) (*PopularMoviesResponse, error) {
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

/*
	=====================
	Первые 1000 популярных фильмов
	=====================
*/

func (c *Client) GetFirst1000PopularMovies() ([]Movie, error) {
	const (
		moviesPerPage = 20
		targetCount   = 1000
		maxPages      = targetCount / moviesPerPage // 50
	)

	result := make([]Movie, 0, targetCount)

	for page := 1; page <= maxPages; page++ {
		log.Println(page)
		resp, err := c.GetPopularMoviesPage(page)
		if err != nil {
			return nil, err
		}

		result = append(result, resp.Results...)
		time.Sleep(250 * time.Millisecond) // защита от rate-limit
	}

	if len(result) > targetCount {
		result = result[:targetCount]
	}

	return result, nil
}

/*
	=====================
	Детали фильма
	=====================
*/

func (c *Client) GetMovieDetails(id int) (map[string]interface{}, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%d?api_key=%s&append_to_response=credits,images,videos",

		id,
		c.apiKey,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb error: status %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
