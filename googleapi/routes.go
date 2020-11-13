package googleapi

import (
	"date-hub-api/server"
)

// GetRoutes ...
func GetRoutes() []server.Route {
	routes := []server.Route{
		server.NewRoute("/", getMyIP, "GET"),
	}
	return routes
}
