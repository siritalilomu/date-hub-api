package googleapi

import (
	"date-hub-api/server"
)

// GetRoutes ...
func GetRoutes() []server.Route {
	routes := []server.Route{
		server.NewRoute("/get-activity", getFood, "GET"),
	}
	return routes
}
