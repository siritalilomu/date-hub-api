package movietventertainment

import (
	"date-hub-api/server"
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
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
}

func getTVShows(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(handleTVShows()); err != nil {
		fmt.Println(err.Error())
	}
}

func handleTVShows() tv {
	var err error
	tvShowURL, tvGenreURL := <-server.DoExternalAPIRequest("GET", fmt.Sprintf("https://api.themoviedb.org/3/discover/tv?api_key=%s&language=en-US&sort_by=popularity.desc&page=1&timezone=America%2FNew_York&include_null_first_air_dates=false", os.Getenv("ApiKey")), "", nil), <-server.DoExternalAPIRequest("GET", fmt.Sprintf("https://api.themoviedb.org/3/genre/tv/list?api_key=%s&language=en-US", os.Getenv("ApiKey")), "", nil)
	defer tvShowURL.Close()

	var t tv

	if err = json.NewDecoder(tvShowURL).Decode(&t); err != nil {
		fmt.Println(err.Error())
	}
	if err = json.NewDecoder(tvGenreURL).Decode(&t); err != nil {
		fmt.Println(err.Error())
	}
	return t
}
