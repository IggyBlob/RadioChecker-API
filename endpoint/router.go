package endpoint

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/IggyBlob/RadioChecker-Core-Library/datastore"
	"errors"
	"github.com/IggyBlob/RadioChecker-Core-Library/metrics"
)

type Config struct {
	DS    *datastore.Datastore
	CORS  string
	Debug bool
	MetricProvider metrics.MetricProvider
}

var conf *Config

// NewRouter creates a new Gorialla mux object based on the given config and assigns routes and loggers.
func NewRouter(c *Config) (*mux.Router, error) {
	if err := validateConfig(c); err != nil {
		return nil, err
	}
	conf = c
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger(handler, route.Name)
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router, nil
}

// validateConfig ensures that the Config object is valid.
func validateConfig(c *Config) error {
	if c == nil {
		return errors.New("config object must not be nil")
	}
	if c.DS == nil {
		return errors.New("config: datastore object must not be nil")
	}
	if c.CORS == "" {
		return errors.New("config: cors string must not be empty")
	}
	if c.MetricProvider == nil {
		return errors.New("config: metric provider must not be nil")
	}

	return nil
}