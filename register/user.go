package register

import (
	"date-hub-api/helpers"
	"date-hub-api/sqlconn"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// User Model ...
type user struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Photo    []byte `json:"photo"`
}

// signup ...
func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup hit")
	var user user
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if user.Password == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 13)
	if err != nil {
		log.Fatal(err)
	}

	user.Password = string(hash)

	sqlconn := sqlconn.Open()
	defer sqlconn.Close()
	if err := sqlconn.Exec("[spcCreateUser] ?, ?, ?, ?", user.Email, user.Name, user.Password, user.Photo); err == 0 {
		http.Error(w, "error on creating user", http.StatusInternalServerError)
		return
	}

	// send back an empty password
	user.Password = ""

	helpers.ResponseJSON(w, user)
}
