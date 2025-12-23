package tmdb

import (
	"encoding/json"
	"fmt"
)

type TVShow struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	FirstAirDate string `json:"first_air_date"`
}

type PopularTVResponse struct {
	Results []TVShow `json:"results"`
}

func (c *Client) GetPopularTVPage(page int) ([]TVShow, error) {
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

	var data PopularTVResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Results, nil
}
