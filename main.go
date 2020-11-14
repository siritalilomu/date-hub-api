package main

import (
	"date-hub-api/googleapi"
	"date-hub-api/movietventertainment"
	"date-hub-api/register"
	"date-hub-api/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
)

type config struct {
	Port string `json:"Port"`
}

// EnvironmentVariables . . .
var EnvironmentVariables config

func main() {

	server := server.NewServer()
	server.AddRoutes(googleapi.GetRoutes())
	server.AddRoutes(register.GetRoutes())
	server.AddRoutes(movietventertainment.GetRoutes())

	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "application/json"})
	origins := handlers.AllowedOrigins([]string{
		"http://localhost:3000",
		"http://localhost:8080",
		"https://localhost:8080",
		"https://localhost:8443",
		"https://date-hub.herokuapp.com/",
		"https://date-hub-backend.herokuapp.com/",
	})

	EnvironmentVariables.Port = os.Getenv("PORT")
	if EnvironmentVariables.Port == "" {
		env, err := os.Open(".env")
		if err != nil {
			fmt.Printf("no config.json file was found: %s\ndefaulting to OS ENV 'PORT'\n", err.Error())
		} else {
			json.NewDecoder(env).Decode(&EnvironmentVariables)
			env.Close()
		}
	}

	if EnvironmentVariables.Port != "" {
		log.Fatal(http.ListenAndServe(":"+EnvironmentVariables.Port, handlers.CORS(methods, origins, headers)(server.Router)))
	} else {
		panic("PORT not set")
	}

}
