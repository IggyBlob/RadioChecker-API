<img src="http://radiochecker.paulhaunschmied.com/assets/img/jack.png" alt="RadioChecker Logo" width="200"/>

# RadioChecker.com API
The API service provides data via a REST interface to other RadioChecker services, e.g. the frontend.

## Language
Go

## Service Configuration
Since the set of available radio stations may be changed at any time (by adding/removing stations), a mature configuration solution for the service is required. [Viper](https://github.com/spf13/viper) is such a library that allows the API service to load config information upon startup from a *config.json* file (watching/re-reading the config file during runtime is not necessary at the moment).

Summarized, the main purposes of the API service's *config.json* are:
+ holding database information/credentials
+ defining the available radio stations, their (REST) URI and database table

Example *config.json*:
```
{
    "service": {
        "name": "rc-api-svc",
        "port": 8080,
        "debug": true,
        "access-control-allow-origin": "*"
    },
    "datastore": {
        "host": "hostname",
        "port": 3306,
        "username": "username",
        "password": "password",
        "schema": "radiochecker"
    }
}
```

## Endpoints (API Design)
The following endpoints are provided by the API service (ordered by importance):

1. `/api/<config.stations.station.uri>/tracks/day/2016-10-01/top`
Most-played tracks (ranked) for the specified day
2. `/api/<config.stations.station.uri>/tracks/day/2016-10-01/all`
Complete list of all tracks played on the specified day
3. `/api/config/stations`
List of available radio stations (those specified in `config.stations`)
4. `/api/<config.stations.station.uri>/tracks/week/52/top`
Most-played tracks (ranked) for the specified week
5. `/api/query/<track name>`
Number of times the track was played on every radio station (how to specify the time period - day/week/year?)
6. `/api/<config.stations.station.uri>/tracks/week/52/all`
Complete list of tracks played on the specified day (not really necessary)


## Deployment
AWS/Bare Metal
