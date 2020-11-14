package register

import (
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Useremail string `json:"useremail"`
		Password  string `json:"password"`
	}

	// handle := func(req request) {}

	// var err error
	// var req request
	// if err = json.NewDecoder(r.Body).Decode(handle(req)); err != nil {
	// 	http.Error(w, "")
	// }
}
