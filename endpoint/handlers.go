package endpoint

import (
	"net/http"
	"fmt"
	"log"
	"github.com/dustin/gojson"
	"github.com/gorilla/mux"
	"time"
	"RadioChecker-Crawler-HitradioOE3/track"
	"strconv"
	"errors"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "RadioChecker API")
}

func getStations(w http.ResponseWriter, r *http.Request) {
	stations, err := ds.GetRadiostations()
	if err != nil {
		log.Printf("getStations Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	j, err := json.Marshal(stations)
	if err != nil {
		log.Printf("getStations Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, j)
}

func getTracksDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, err := ds.GetRadiostationID(vars["station"]); err != nil {
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

	var tracks []track.Track
	if vars["filter"] == "top" {
		tracks, err = ds.GetTopTracks(vars["station"], since, until)
	} else {
		tracks, err = ds.GetAllTracks(vars["station"], since, until)
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
		Plays []track.Track `json:"plays"`
	}

	resp := response{vars["station"], vars["date"], tracks}
	j, err := json.Marshal(resp)
	if err != nil {
		log.Printf("getTracksDay Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, j)
}

func getTracksWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, err := ds.GetRadiostationID(vars["station"]); err != nil {
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
	
	log.Println(since)
	log.Println(until)

	var tracks []track.Track
	if vars["filter"] == "top" {
		tracks, err = ds.GetTopTracks(vars["station"], since, until)
	} else {
		tracks, err = ds.GetAllTracks(vars["station"], since, until)
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
		Plays []track.Track `json:"plays"`
	}

	resp := response{vars["station"], vars["date"], tracks}
	j, err := json.Marshal(resp)
	if err != nil {
		log.Printf("getTracksDay Handler: %s\n", err.Error())
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, j)
}

func getTrackQuery(w http.ResponseWriter, r *http.Request) {
	handleError(w, http.StatusNotImplemented, "Not implemented")
}

func writeJSONResponse(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func handleError(w http.ResponseWriter, statuscode int, msg string) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(statuscode)
	if msg != "" {
		w.Write([]byte(msg + "\n"))
	}
}

func firstDayOfISOWeek(year int, week int, timezone *time.Location) (time.Time, error) {
	if week < 1 || week > 53 {
		return time.Time{}, errors.New("week out of range")
	}

	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date, nil
}
