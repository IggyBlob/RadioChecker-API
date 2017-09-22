package endpoint

import (
	"net/http"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"time"
	"github.com/IggyBlob/RadioChecker-Core-Library/model"
	"strconv"
	"strings"
)

// index is the default handler.
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "RadioChecker API v1.0\n\nCopyright (C) 2017 Paul Haunschmied.\nwww.radiochecker.com")
}

// getStations returns a map of all active radiostations using the format { "Name":"URI", ... }.
func getStations(w http.ResponseWriter, r *http.Request) {
	stations, err := conf.DS.GetRadiostations()
	if err != nil {
		log.Printf("getStations Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	j, err := jsonMarshal(stations)
	if err != nil {
		log.Printf("getStations Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, j)
}

// getTracksDay returns either the top-3 tracks or all tracks (without duplicates) of a day.
func getTracksDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, err := conf.DS.GetRadiostationID(vars["station"]); err != nil {
		log.Printf("getTracksDay Handler: GetRadiostationID(%s): %s\n", vars["station"], err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	t, err := time.Parse("2006-01-02", vars["date"])
	if err !=  nil {
		log.Printf("getTracksDay Handler: time.Parse(%s): %s\n", vars["date"], err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	loc, _ := time.LoadLocation("Europe/Vienna")
	since := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0,
		loc,
	)
	until := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		23, 59, 59, 0,
		loc,
	)

	var tracks []model.Track
	if vars["filter"] == "top" {
		tracks, err = conf.DS.GetTopTracks(vars["station"], since, until)
	} else {
		tracks, err = conf.DS.GetAllTracks(vars["station"], since, until)
	}
	if err != nil {
		log.Printf("getTracksDay Handler: GetTopTracks/GetAllTracks(%s, %q, %q): %s\n", vars["station"],
			since, until, err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	type response struct {
		Station string `json:"station"`
		Date string `json:"date"`
		Plays []model.Track `json:"plays"`
	}

	resp := response{vars["station"], vars["date"], tracks}
	j, err := jsonMarshal(resp)
	if err != nil {
		log.Printf("getTracksDay Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, j)
}

// getTracksWeek returns either the top-3 tracks or all tracks (without duplicates) of a day.
func getTracksWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, err := conf.DS.GetRadiostationID(vars["station"]); err != nil {
		log.Printf("getTracksWeek Handler: GetRadiostationID(%s): %s\n", vars["station"], err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	loc, _ := time.LoadLocation("Europe/Vienna")

	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		log.Printf("getTracksWeek Handler: atoi(%s): %s\n", vars["year"], err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}
	week, err := strconv.Atoi(vars["week"])
	if err != nil {
		log.Printf("getTracksWeek Handler: atoi(%s): %s\n", vars["week"], err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	since, err := firstDayOfISOWeek(year, week, loc) // Monday 00:00:00
	if err != nil {
		log.Printf("getTracksWeek Handler: firstDayOfISOWeek(%d, %d, %q): %s\n", year, week, loc, err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	until := since.AddDate(0, 0, 6) // Sunday 00:00:00
	until = time.Date(
		until.Year(),
		until.Month(),
		until.Day(),
		23, 59, 59, 0,
		since.Location(),
	) // Sunday 23:59:59

	var tracks []model.Track
	if vars["filter"] == "top" {
		tracks, err = conf.DS.GetTopTracks(vars["station"], since, until)
	} else {
		tracks, err = conf.DS.GetAllTracks(vars["station"], since, until)
	}
	if err != nil {
		log.Printf("getTracksWeek Handler: GetTopTracks/GetAllTracks(%s, %q, %q): %s\n", vars["station"],
			since, until, err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	type response struct {
		Station string `json:"station"`
		WeekNo string `json:"weekNo"`
		BeginDate string `json:"beginDate"`
		EndDate string `json:"endDate"`
		Plays []model.Track `json:"plays"`
	}

	resp := response{
		vars["station"],
		vars["week"],
		since.Format("2006-01-02"),
		until.Format("2006-01-02"),
		tracks,
	}
	j, err := jsonMarshal(resp)
	if err != nil {
		log.Printf("getTracksWeek Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, j)
}

// getSearchQueryDay returns the times a track has been played on the specified day on every active radiostation.
func getSearchQueryDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := time.Parse("2006-01-02", vars["date"])
	if err !=  nil {
		log.Printf("getSearchQueryDay Handler: time.Parse(%s): %s\n", vars["date"], err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	loc, _ := time.LoadLocation("Europe/Vienna")
	since := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0,
		loc,
	)
	until := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		23, 59, 59, 0,
		loc,
	)

	tracks, err := conf.DS.GetSearchResult(strings.Replace(vars["query"], "+", ",", -1), since, until)
	if err != nil {
		log.Printf("getSearchQueryDay Handler: GetSearchResult(%s, %q, %q): %s\n",
			strings.Replace(vars["query"], "+", ",", -1), since, until, err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	resp, err := orderSearchResults(tracks, since, time.Time{}, 0)
	if err != nil {
		log.Printf("getSearchQueryWeek Handler: orderSearchResults(%q, %q, %q, %d): %s\n",
			tracks, since, time.Time{}, 0, err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	j, err := jsonMarshal(resp)
	if err != nil {
		log.Printf("getSearchQueryWeek Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, j)
}

// getSearchQueryWeek returns the times a track has been played during the specified week on every active radiostation.
func getSearchQueryWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	loc, _ := time.LoadLocation("Europe/Vienna")

	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		log.Printf("getSearchQueryWeek Handler: atoi(%s): %s\n", vars["year"], err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}
	week, err := strconv.Atoi(vars["week"])
	if err != nil {
		log.Printf("getSearchQueryWeek Handler: atoi(%s): %s\n", vars["week"], err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	since, err := firstDayOfISOWeek(year, week, loc) // Monday 00:00:00
	if err != nil {
		log.Printf("getSearchQueryWeek Handler: firstDayOfISOWeek(%d, %d, %q): %s\n", year, week, loc,
			err.Error())
		handleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	until := since.AddDate(0, 0, 6) // Sunday 00:00:00
	until = time.Date(
		until.Year(),
		until.Month(),
		until.Day(),
		23, 59, 59, 0,
		since.Location(),
	) // Sunday 23:59:59

	tracks, err := conf.DS.GetSearchResult(strings.Replace(vars["query"], "+", ",", -1), since, until)
	if err != nil {
		log.Printf("getSearchQueryWeek Handler: GetSearchResult(%s, %q, %q): %s\n",
			strings.Replace(vars["query"], "+", ",", -1), since, until, err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	resp, err := orderSearchResults(tracks, since, until, week)
	if err != nil {
		log.Printf("getSearchQueryWeek Handler: orderSearchResults(%q, %q, %q, %d): %s\n",
			tracks, since, until, week, err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	j, err := jsonMarshal(resp)
	if err != nil {
		log.Printf("getSearchQueryWeek Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, j)
}

// getSearchQueryYear returns the times a track has been played during the specified year on every active radiostation.
func getSearchQueryYear(w http.ResponseWriter, r *http.Request) {
	handleNotImplemented(w)
}
