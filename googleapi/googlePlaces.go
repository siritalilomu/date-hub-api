package vcaas

import (
	"date-hub/server"
)

// GetRoutes ...
func GetRoutes() []server.Route {
	routes := []server.Route{
		server.NewRoute("/", events, "POST"),
	}
	return routes
}
