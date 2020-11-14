package movietventertainment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type tv struct {
	Results []struct {
		Name          string   `json:"name"`
		GenreIDs      []int    `json:"genre_ids"`
		VoteCount     int      `json:"vote_count"`
		ID            int      `json:"id"`
		OverView      string   `json:"overview"`
		OriginCountry []string `json:"origin_country"`
	} `json:"results"`
}

func getTVShows(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(handleTVShows()); err != nil {
		fmt.Println(err.Error())
	}
}

func handleTVShows() tv {
	tvShowURL := fmt.Sprintf("https://api.themoviedb.org/3/discover/tv?api_key=%s&language=en-US&sort_by=popularity.desc&page=1&timezone=America%2FNew_York&include_null_first_air_dates=false", os.Getenv("ApiKey"))
	tvGenreURL := fmt.Sprintf("https://api.themoviedb.org/3/genre/tv/list?api_key=%s&language=en-US", os.Getenv("ApiKey"))

	tvShowReq, err := http.Get(tvShowURL)
	if err != nil {
		fmt.Println(err.Error())
	}

	var t tv

	if err = json.NewDecoder(tvShowReq.Body).Decode(&t); err != nil {
		fmt.Println(err.Error())
	}
	return t
}
