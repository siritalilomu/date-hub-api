package movietventertainment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type movies struct {
	Results []struct {
		Popularity  float32 `json:"popularity"`
		VoteCount   int32   `json:"vote_count"`
		Video       bool    `json:"video"`
		ID          int     `json:"id"`
		GenreIDs    []int   `json:"genre_ids"`
		Title       string  `json:"title"`
		VoteAverage float32 `json:"vote_average"`
		Overview    string  `json:"overview"`
		ReleaseDate string  `json:"release_date"`
	} `json:"results"`
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(getMovieHandler()); err != nil {
		fmt.Println(err.Error())
	}
}

func getMovieHandler() movies {

	url := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?api_key=%s&language=en-US&sort_by=popularity.desc&include_adult=false&include_video=false&page=1", os.Getenv("ApiKey"))
	fmt.Println(url)
	req, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer req.Body.Close()
	var m movies
	if err = json.NewDecoder(req.Body).Decode(&m); err != nil {
		fmt.Println(err.Error())
	}
	return m
}
