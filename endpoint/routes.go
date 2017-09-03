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
		"GetSearchQueryDay",
		"GET",
		"/search/day/{date:\\d{4}-\\d{2}-\\d{2}}/{query:\\w[+\\w]+}",
		getSearchQueryDay,
	},
	Route{
		"GetSearchQueryWeek",
		"GET",
		"/search/week/{year:\\d{4}}/{week:\\d{2}}/{query:\\w[+\\w]+}",
		getSearchQueryWeek,
	},
	Route{
		"GetSearchQueryYear",
		"GET",
		"/search/year/{year:\\d{4}}/{query:\\w[+\\w]+}",
		getSearchQueryYear,
	},
}
