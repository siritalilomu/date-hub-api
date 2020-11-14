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

// login ...
// func login(w http.ResponseWriter, r *http.Request) {

// 		var user user
// 		error := &serverError{}

// 		json.NewDecoder(r.Body).Decode(&user)

// 		if user.Email == "" {
// 			error.Message = "Email is required"
// 			helpers.RespondWithError(w, http.StatusBadRequest, error.Message)
// 			return
// 		}

// 		if user.Password == "" {
// 			error.Message = "Password is required"
// 			helpers.RespondWithError(w, http.StatusBadRequest, error.Message)
// 			return
// 		}

// 		password := user.Password

// 		sqlconn := sqlconn.Open()
// 		defer sqlconn.Close()
// 		row := sqlconn.QueryRow("[spcCreateUser] ?", user.Email)
// 		err := row.Scan(&user.Email, &user.Password)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				error.Message = "Invalid User and Password"
// 				helpers.RespondWithError(w, http.StatusBadRequest, error.Message)
// 				return
// 			}
// 			log.Fatal(err)
// 		}

// 		// Compare password with hashed password
// 		hashedPassword := user.Password
// 		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// 		if err != nil {
// 			error.Message = "Invalid User and Password"
// 			helpers.RespondWithError(w, http.StatusUnauthorized, error.Message)
// 			return
// 		}

// 		token, err := helpers.GenerateToken(user.Email)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		jwt.Token = token

// 		helpers.ResponseJSON(w, jwt)
// 	}
// }
