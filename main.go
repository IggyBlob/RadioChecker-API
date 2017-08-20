package main

import (
	"github.com/spf13/viper"
	"fmt"
	"RadioChecker-API/endpoint"
	"log"
	"net/http"
)

const confFile string = "config_prod"
const confDir string = "."

func main() {
	viper.SetConfigName(confFile)
	viper.AddConfigPath(confDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error on reading config file: %s \n", err))
	}

	router := endpoint.NewRouter()
	log.Fatal(http.ListenAndServe(":" + viper.GetString("service.port"), router))
}
