package anilist

type PopularAnimeResponse struct {
	Data struct {
		Page struct {
			Media []AnimeDTO `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

type AnimeDTO struct {
	ID    int  `json:"id"`
	IDMal *int `json:"idMal"`

	Title struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
		Native  string `json:"native"`
	} `json:"title"`

	Description string `json:"description"`

	Format     string `json:"format"`
	Status     string `json:"status"`
	Season     string `json:"season"`
	SeasonYear *int   `json:"seasonYear"`

	Episodes *int `json:"episodes"`
	Duration *int `json:"duration"`

	StartDate *DateDTO `json:"startDate"`
	EndDate   *DateDTO `json:"endDate"`

	Genres []string `json:"genres"`

	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`

	Studios struct {
		Nodes []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"nodes"`
	} `json:"studios"`

	Characters struct {
		Nodes []struct {
			Name struct {
				Full string `json:"full"`
			} `json:"name"`
		} `json:"nodes"`
	} `json:"characters"`

	VoiceActors struct {
		Nodes []struct {
			Name struct {
				Full string `json:"full"`
			} `json:"name"`
		} `json:"nodes"`
	} `json:"voiceActors"`

	MeanScore  *int `json:"meanScore"`
	Popularity int  `json:"popularity"`
	Favourites int  `json:"favourites"`

	CoverImage struct {
		Large string `json:"large"`
	} `json:"coverImage"`

	Trailer *struct {
		Site string `json:"site"`
		ID   string `json:"id"`
	} `json:"trailer"`
}

type DateDTO struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}
