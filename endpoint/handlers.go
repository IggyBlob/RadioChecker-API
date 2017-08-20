package endpoint

import (
	"net/http"
	"fmt"
	"log"
	"github.com/dustin/gojson"
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

func getTopTracks(w http.ResponseWriter, r *http.Request) {
	handleError(w, http.StatusNotImplemented, "Not implemented")

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
