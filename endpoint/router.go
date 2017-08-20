package endpoint

import (
	"github.com/gorilla/mux"
	"net/http"
	"RadioChecker-Crawler-HitradioOE3/datastore"
	"errors"
)

var ds *datastore.Datastore

func NewRouter(d *datastore.Datastore) (*mux.Router, error) {
	if d == nil {
		return nil, errors.New("datastore object must not be nil")
	}
	ds = d
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/api/")
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger(handler, route.Name)
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router, nil
}
