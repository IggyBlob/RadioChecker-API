package main

import (
	"github.com/spf13/viper"
	"RadioChecker-API/endpoint"
	"log"
	"net/http"
	"RadioChecker-Crawler-HitradioOE3/datastore"
)

const confFile string = "config_prod"
const confDir string = "."

func main() {
	viper.SetConfigName(confFile)
	viper.AddConfigPath(confDir)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error on reading config file: %s \n", err)
	}

	ds, err := datastore.NewDatastore(
		viper.GetString("datastore.username"),
		viper.GetString("datastore.password"),
		viper.GetString("datastore.host"),
		viper.GetInt("datastore.port"),
		viper.GetString("datastore.schema"),
	)
	if err != nil {
		log.Fatalf("Fatal error on creating datastore: %s \n", err)
	}
	defer ds.Close()

	router, err := endpoint.NewRouter(ds)
	if err != nil {
		log.Fatalf("Unable to create router: %s\n", err)
	}
	log.Printf("listening")
	log.Fatal(http.ListenAndServe(":" + viper.GetString("service.port"), router))
}
