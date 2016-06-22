package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	forecast "github.com/jtslear/forecast/v2"
)

func check_config(environmentVariable string) {
	if environmentVariable == "" {
		log.Fatal("An environment variable is missing.  Exiting...")
	}
}

func main() {
	api_key := os.Getenv("FORCAST_IO_API_KEY")
	debug := os.Getenv("DEBUG")
	influxdb_connection_string := os.Getenv("INFLUXDB_CONNECTION_STRING")
	influxdb_database := os.Getenv("INFLUXDB_DATABASE")
	influxdb_password := os.Getenv("INFLUXDB_PASSWORD")
	influxdb_username := os.Getenv("INFLUXDB_USERNAME")
	latitude := os.Getenv("LATITUDE")
	longitude := os.Getenv("LONGITUDE")

	check_config(api_key)
	check_config(latitude)
	check_config(longitude)
	check_config(influxdb_connection_string)
	check_config(influxdb_database)
	check_config(influxdb_password)
	check_config(influxdb_username)

	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxdb_connection_string,
		Username: influxdb_username,
		Password: influxdb_password,
	})

	if err != nil {
		log.Fatal("Unable to create InfluxDB Connection: ", err, "Exiting...")
	}

	influxdbPingTime, _, err := influxdbClient.Ping(time.Duration(10))

	if err != nil {
		log.Fatal("Unable to maintain InfluxDB Connection: ", err, "Exiting...")
	}
	if debug == "true" {
		fmt.Println("InfluxDB is alive, response: ", influxdbPingTime)
	}

	batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxdb_database,
		Precision: "s",
	})

	if err != nil {
		log.Fatalln("Unable to create a batch configuration: ", err)
	}

	f, err := forecast.Get(api_key, latitude, longitude, "now", forecast.US)
	if err != nil {
		log.Fatal(err)
	}
	time := time.Now()
	humidity := (f.Currently.Humidity * 100)
	dataPoints := map[string]interface{}{
		"humidity":    humidity,
		"temperature": f.Currently.Temperature,
		"pressure":    f.Currently.Pressure,
		"dew_point":   f.Currently.DewPoint,
	}
	if debug == "true" {
		fmt.Printf("humidity: %.2f%%\n", humidity)
		fmt.Printf("temperature: %.2f F\n", f.Currently.Temperature)
		fmt.Printf("pressure: %.2f millibars\n", f.Currently.Pressure)
		fmt.Printf("dew_point: %.2f F\n", f.Currently.DewPoint)
	}
	data, err := client.NewPoint("weather", nil, dataPoints, time)
	if err != nil {
		log.Fatalln("Unable to create NewPoint: ", err)
	}
	batchPoints.AddPoint(data)

	influxdbClient.Write(batchPoints)
	if err != nil {
		log.Fatal("Unable to write data to influxdb: ", err)
	}
}
