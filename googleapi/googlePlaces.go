package googleapi

import (
	"date-hub-api/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func getPhoto(ref string) (img string) {

	URL := fmt.Sprintf(`https://maps.googleapis.com/maps/api/place/photo?maxwidth=1000&photoreference=%s&key=%s`, ref, os.Getenv("GOOGLE_KEY"))
	r, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		server.PanicWithStatus(err, http.StatusBadRequest)
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		server.PanicWithStatus(err, http.StatusBadRequest)
	}
	defer resp.Body.Close()

	var resbody []byte
	resbody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	rb := fmt.Sprintf(`%s`, string(resbody))

	return rb
}

func getFood(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Lat     string
		Lon     string
		Type    string
		Keyword string
	}

	type response struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
			Name   string `json:"name"`
			Photos []struct {
				Height           int      `json:"height"`
				HTMLAttributions []string `json:"html_attributions"`
				PhotoReference   string   `json:"photo_reference"`
				ImageUrl         string   `json:"imgUrl"`
				Width            int      `json:"width"`
			} `json:"photos"`
			PlaceID          string  `json:"place_id"`
			Rating           float64 `json:"rating"`
			Reference        string  `json:"reference"`
			UserRatingsTotal int     `json:"user_ratings_total"`
			Vicinity         string  `json:"vicinity"`
		} `json:"results"`
	}

	handler := func(req request) *response {
		URL := fmt.Sprintf(`https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s,%s&radius=8000&type=%s&keyword=%s&key=%s`, req.Lat, req.Lon, req.Type, req.Keyword, os.Getenv("GOOGLE_KEY"))

		r, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			server.PanicWithStatus(err, http.StatusBadRequest)
		}

		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			server.PanicWithStatus(err, http.StatusBadRequest)
		}
		defer resp.Body.Close()

		var res response
		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			log.Println(err)
		}

		if len(res.Results[0].Photos[0].PhotoReference) > 10 {
			res.Results[0].Photos[0].ImageUrl = getPhoto(res.Results[0].Photos[0].PhotoReference)
		}

		if resp.StatusCode > 204 {
			panic(fmt.Errorf("expected status code 200 or 204 but got %d", resp.StatusCode))
		}

		return &res
	}

	var req request = request{Lat: server.GetStringParam(r, "lat"), Lon: server.GetStringParam(r, "lon"), Type: server.GetStringParam(r, "type"), Keyword: server.GetStringParam(r, "keyword", true)}

	res := handler(req)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func getActivity(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Lat  string
		Lon  string
		Type []string
	}

	type responsebody struct {
		Results []struct {
			BusinessStatus string `json:"business_status"`
			Geometry       struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
			Name   string `json:"name"`
			Photos []struct {
				Height           int      `json:"height"`
				HTMLAttributions []string `json:"html_attributions"`
				PhotoReference   string   `json:"photo_reference"`
				Width            int      `json:"width"`
			} `json:"photos"`
			PlaceID          string  `json:"place_id"`
			Rating           float64 `json:"rating"`
			UserRatingsTotal int     `json:"user_ratings_total"`
			Vicinity         string  `json:"vicinity"`
		} `json:"results"`
	}

	type httpresponse struct {
		url      string
		response *http.Response
		err      error
		res      responsebody
	}

	handler := func(req request) []*responsebody {
		fmt.Println(req)
		var urls []string
		for _, t := range req.Type {
			urls = append(urls, fmt.Sprintf(`https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s,%s&radius=8000&type=%s&key=%s`, req.Lat, req.Lon, t, os.Getenv("GOOGLE_KEY")))
		}
		var responses []*responsebody
		fmt.Println(urls)
		var wg sync.WaitGroup
		for _, url := range urls {
			wg.Add(1)
			go func(url string) {
				r, err := http.NewRequest("GET", url, nil)
				if err != nil {
					server.PanicWithStatus(err, http.StatusBadRequest)
				}
				client := &http.Client{}
				resp, err := client.Do(r)
				if err != nil {
					server.PanicWithStatus(err, http.StatusBadRequest)
				}
				defer resp.Body.Close()

				var resbody []byte
				resbody, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				rb := &responsebody{}
				err = json.Unmarshal(resbody, rb)
				if err != nil {
					panic(err)
				}
				responses = append(responses, rb)
				wg.Done()
			}(url)
		}
		wg.Wait()

		return responses
	}

	var req request = request{Lat: server.GetStringParam(r, "lat"), Lon: server.GetStringParam(r, "lon"), Type: server.GetStringParams(r, "type")}

	res := handler(req)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

// amusement_park
// aquarium
// art_gallery
// beauty_salon
// bowling_alley
// campground
// casino
// gym
// movie_theater
// museum
// night_club
// park
// spa
// zoo
