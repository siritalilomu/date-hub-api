package dates

import (
	"date-hub-api/helpers"
	"date-hub-api/sqlconn"
	"encoding/json"
	"net/http"
)

func createDate(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserID   int     `json:"userId"`
		DateDate string  `json:"dateDate"`
		Lat      float32 `json:"latitude"`
		Lon      float32 `json:"longitude"`
	}

	type response struct {
		DateID int `json:"dateId"`
	}

	handler := func(req request) response {

		db := sqlconn.Open()
		defer db.Close()
		var resp response
		if err := db.QueryRow("[spcCreateDate] ?, ?, ?, ?", req.UserID, req.DateDate, req.Lat, req.Lon).Scan(&resp.DateID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return resp
	}
	var err error
	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	helpers.ResponseJSON(w, handler(req))
}
