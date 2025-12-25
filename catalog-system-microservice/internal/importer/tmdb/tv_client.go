package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
=====================
DTO для сериалов
=====================
*/

type Series struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	FirstAirDate string  `json:"first_air_date"`
	Popularity   float64 `json:"popularity"`
}

type PopularSeriesResponse struct {
	Page       int      `json:"page"`
	Results    []Series `json:"results"`
	TotalPages int      `json:"total_pages"`
}

func (c *Client) GetPopularSeriesPage(page int) (*PopularSeriesResponse, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/tv/popular?api_key=%s&page=%d",
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

	var data PopularSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) GetSeriesDetails(id int) (map[string]interface{}, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/tv/%d?api_key=%s&append_to_response=credits,images,videos",
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
