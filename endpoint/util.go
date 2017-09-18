package endpoint

import (
	"net/http"
	"time"
	"errors"
	"github.com/IggyBlob/RadioChecker-Core-Library/model"
)

type orderSearchResult struct {
	Track model.Track `json:"track"`
	Plays []int `json:"plays"`
}

type getSearchQueryResponse struct {
	Stations []string `json:"stations"`
	Date string `json:"date,omitempty"`
	WeekNo int `json:"weekNo,omitempty"`
	BeginDate string `json:"beginDate,omitempty"`
	EndDate string `json:"endDate,omitempty"`
	Results []*orderSearchResult`json:"results"`
}

// writeJSONResponse is a utility function that writes a 200 OK JSON response to the ResponseWriter.
func writeJSONResponse(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// handleError is a utility function that writes a specified error response to the ResponseWriter.
func handleError(w http.ResponseWriter, statuscode int, msg string) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statuscode)
	if msg != "" {
		w.Write([]byte(msg + "\n"))
	}
}

// handleNotImplemented is a utility function that writes a 501 Not implemented to the ResponseWriter.
func handleNotImplemented(w http.ResponseWriter) {
	handleError(w, http.StatusNotImplemented, "Not implemented")
}

// firstDayOfISOWeek is a utility function that returns the first date of a specified week.
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

// orderSearchResults is a utility function that takes the database search result and converts it into a result object
// that can be marshalled into JSON.
func orderSearchResults(searchResult []model.Track, since time.Time, until time.Time,
weekNo int) (*getSearchQueryResponse, error){

	if searchResult == nil || since.IsZero() {
		return nil, errors.New("searchResult and/or since must not be empty")
	}

	if (until.IsZero() && weekNo != 0) || (!until.IsZero() && (weekNo < 1 || weekNo > 53)) {
		return nil, errors.New("until and weekNo must either be nil-valued or within a valid rang")
	}

	// group equal radio stations by storing their Name value into the stationsFiltered array
	// create a stationsMeta map that assigns the position of the radiostation's name in the stationsFiltered array
	// to make inserting the radiostation's play count to the correct position of a return object's Play array
	// easier (see group equal tracks)
	stationsMeta := make(map[string]int)
	stationsFiltered := make([]string, 0)
	i := 0
	for _, track := range searchResult {
		if _, exists := stationsMeta[track.Radiostation.URI]; !exists {
			stationsMeta[track.Radiostation.URI] = i
			stationsFiltered = append(stationsFiltered, track.Radiostation.Name)
			i++
		}
	}

	// group equal tracks in a new result object and assign their play count to the respective position in the
	// result object's Play array. Before that, persist the absolute position of the track in the result set into
	// trackOrder to be able to restore the relevance order determined by the database engine.
	trackOrder := make(map[int64]int)
	tracksGrouped := make(map[int64]*orderSearchResult)
	i = 0
	for _, track := range searchResult {
		if _, exists := tracksGrouped[track.ID]; !exists {
			trackOrder[track.ID] = i
			tracksGrouped[track.ID] = &orderSearchResult{ track, make([]int, len(stationsFiltered))}
			i++
		}
		tracksGrouped[track.ID].Plays[stationsMeta[track.Radiostation.URI]] = track.Count
		tracksGrouped[track.ID].Track.Count = 0 // prevent JSON marshalling of the original Count field by
							// setting it to its NULL value
	}

	// convert map into array using the relevance order determined by the database engine
	results := make([]*orderSearchResult, len(tracksGrouped))
	for _, result := range tracksGrouped {
		results[trackOrder[result.Track.ID]] = result
	}

	// build getSearchQueryResponse object for either a day or a whole week, depending on whether until is set or
	// not
	resp := &getSearchQueryResponse{ Stations: stationsFiltered, Results: results }
	if since.Before(until) {
		resp.BeginDate = since.Format("2006-01-02")
		resp.EndDate = until.Format("2006-01-02")
		resp.WeekNo = weekNo
	} else {
		resp.Date = since.Format("2006-01-02")
	}

	return resp, nil
}
