package register

import (
	"date-hub-api/helpers"
	"date-hub-api/server"
	"date-hub-api/sqlconn"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type request struct {
	Useremail string `json:"useremail"`
	Password  string `json:"password"`
}

type response struct {
	UserName    string
	UserEmail   string
	UserID      int
	PhotoBinary string
	DateID      *int
	Rating      *int
	Dates       *[]string
}

func login(w http.ResponseWriter, r *http.Request) {

	resp, statusCode := handleLogin(r)

	if statusCode == http.StatusUnauthorized {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		helpers.ResponseJSON(w, resp)
	}

}

func handleLogin(r *http.Request) (response, int) {
	var err error

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.PanicWithStatus(err, http.StatusInternalServerError)
	}
	db := sqlconn.Open()
	defer db.Close()

	var passwordHash string

	if err = db.QueryRow("[spcGetUserHash] ?", req.Useremail).Scan(&passwordHash); err != nil {
		server.PanicWithStatus(err, http.StatusInternalServerError)
	}
	var u response
	if err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return u, http.StatusUnauthorized
	}

	if err = db.QueryRow("[spcGetUserData_login] ?", req.Useremail).Scan(
		&u.UserName,
		&u.UserEmail,
		&u.UserID,
		&u.PhotoBinary,
		&u.DateID,
		&u.Rating,
		&u.Dates); err != nil {
		log.Fatal(err)
	}

	return u, 200
}
