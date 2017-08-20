package endpoint

import (
	"net/http"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes {
	Route{
		"Index",
		"GET",
		"/",
		index,
	},
	Route{
		"GetStations",
		"GET",
		"/stations",
		getStations,
	},
	Route{
		"GetTopTracks",
		"GET",
		"/{station}/tracks/{period:}/{date:dddd-dd-dd|dd}/{filter}",
		getTopTracks,
	},
	Route{
		"GetTrackQuery",
		"GET",
		"/query/{trackname}",
		getTrackQuery,
	},
}
