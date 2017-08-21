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
		"GetTracksDay",
		"GET",
		"/{station}/tracks/day/{date:\\d{4}-\\d{2}-\\d{2}}/{filter:all|top}",
		getTracksDay,
	},
	Route{
		"GetTracksWeek",
		"GET",
		"/{station}/tracks/week/{year:\\d{4}}/{week:\\d{2}}/{filter:all|top}",
		getTracksWeek,
	},
	Route{
		"GetTrackQuery",
		"GET",
		"/query/{artist}/{title}/day/{date}",
		getTrackQueryDay,
	},
	Route{
		"GetTrackQuery",
		"GET",
		"/query/{artist}/{title}/week/{week}",
		getTrackQueryWeek,
	},
	Route{
		"GetTrackQuery",
		"GET",
		"/query/{artist}/{title}/year/{year}",
		getTrackQueryYear,
	},
}
