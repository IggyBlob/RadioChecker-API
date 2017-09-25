package main

import (
	"github.com/spf13/viper"
	"RadioChecker-API/endpoint"
	"log"
	"net/http"
	"github.com/IggyBlob/RadioChecker-Core-Library/datastore"
	"github.com/IggyBlob/RadioChecker-Core-Library/metrics"
)

const confFile string = "config_prod"
const confDir string = "."

func main() {
	viper.SetConfigName(confFile)
	viper.AddConfigPath(confDir)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error on reading config file: %s\n", err)
	}

	ds, err := datastore.NewDatastore(
		viper.GetString("datastore.username"),
		viper.GetString("datastore.password"),
		viper.GetString("datastore.host"),
		viper.GetInt("datastore.port"),
		viper.GetString("datastore.schema"),
	)
	if err != nil {
		log.Fatalf("Fatal error on creating datastore: %s\n", err)
	}
	defer ds.Close()

	amp, err := metrics.NewAWSMetricProvider(
		viper.GetString("metrics.provider.awsAkid"),
		viper.GetString("metrics.provider.awsSecretKey"),
		viper.GetString("metrics.provider.awsRegion"),
		viper.GetBool("metrics.enabled"),
	)
	if err != nil {
		log.Fatalf("Fatal error on creating AWSMetricProvider: %s\n", err)
	}

	config := &endpoint.Config{
		DS: ds,
		CORS: viper.GetString("service.access-control-allow-origin"),
		Debug: viper.GetBool("service.debug"),
		MetricProvider: amp,
	}

	router, err := endpoint.NewRouter(config)
	if err != nil {
		log.Fatalf("Unable to create router: %s\n", err)
	}
	log.Printf("%s up and running on port %d", viper.GetString("service.name"), viper.GetInt("service.port"))
	log.Fatal(http.ListenAndServe(":" + viper.GetString("service.port"), router))
}
