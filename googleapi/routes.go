package googleapi

import (
	"date-hub-api/server"
)

// GetRoutes ...
func GetRoutes() []server.Route {
	routes := []server.Route{
		server.NewRoute("/get-food", getFood, "GET"),
		server.NewRoute("/get-activities", getActivity, "GET"),
	}
	return routes
}
