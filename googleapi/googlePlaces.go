package googleapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func getMyIP(w http.ResponseWriter, r *http.Request) {

	type response struct {
		IPAddress string `json:"IPAddress"`
	}

	handler := func() *response {

		var r *http.Response
		var err error
		if r, err = http.Get("http://icanhazip.com"); err != nil {
			panic(err)
		}

		var res response
		var ip []byte
		ip, _ = ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		res.IPAddress = string(ip)

		return &res
	}

	resp := handler()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
