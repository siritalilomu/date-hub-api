package movietventertainment

import "date-hub-api/server"

// GetRoutes ...
func GetRoutes() []server.Route {
	return []server.Route{server.NewRoute("/theaters", getTheaters, "GET"),
		server.NewRoute("/movies", getMovie, "GET"),
		server.NewRoute("/tv/shows", getTVShows, "GET")}
}
