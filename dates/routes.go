package dates

import "date-hub-api/server"

// GetRoutes . . .
func GetRoutes() []server.Route {
	return []server.Route{server.NewRoute("/new-date", createDate, "POST")}
}
