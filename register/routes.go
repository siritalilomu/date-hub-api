package register

import (
	"date-hub-api/server"
)

// GetRoutes ...
func GetRoutes() []server.Route {

	routes := []server.Route{
		server.NewRoute("/signup", signup, "POST"),
		// server.NewRoute("/login", login, "POST"),
	}
	return routes
}
