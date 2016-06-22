package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/jsgoecke/nest"
)

func check_config(environmentVariable string) {
	if environmentVariable == "" {
		log.Fatal("An environment variable is missing.  Exiting...")
	}
}

func main() {
	debug := os.Getenv("DEBUG")
	influxdb_connection_string := os.Getenv("INFLUXDB_CONNECTION_STRING")
	influxdb_database := os.Getenv("INFLUXDB_DATABASE")
	influxdb_password := os.Getenv("INFLUXDB_PASSWORD")
	influxdb_username := os.Getenv("INFLUXDB_USERNAME")
	clientID := os.Getenv("NEST_CLIENT_ID")
	state := os.Getenv("STATE")
	clientSecret := os.Getenv("NEST_CLIENT_SECRET")
	authorizationCode := os.Getenv("NEST_AUTH_CODE")
	token := os.Getenv("NEST_TOKEN")

	check_config(clientID)
	check_config(clientSecret)
	check_config(authorizationCode)
	check_config(token)
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

	time := time.Now()
	nestClient := nest.New(clientID, state, clientSecret, authorizationCode)
	nestClient.Token = token
	devices, _ := nestClient.Devices()
	if len(devices.Thermostats) == 0 {
		log.Fatalln("Had some trouble finding devices.")
	}
	for device := range devices.Thermostats {
		dataPoints := map[string]interface{}{
			"name":        devices.Thermostats[device].Name,
			"online":      devices.Thermostats[device].IsOnline,
			"has_leaf":    devices.Thermostats[device].HasLeaf,
			"hvac_mode":   devices.Thermostats[device].HvacMode,
			"temperature": devices.Thermostats[device].AmbientTemperatureF,
			"humidify":    devices.Thermostats[device].Humidity,
			"hvac_state":  devices.Thermostats[device].HvacState,
		}
		if debug == "true" {
			fmt.Printf("Device ID:\t%v\n", devices.Thermostats[device].DeviceID)
			fmt.Printf("Name:\t\t%v\n", devices.Thermostats[device].Name)
			fmt.Printf("Online?:\t%v\n", devices.Thermostats[device].IsOnline)
			fmt.Printf("Has Leaf:\t%v\n", devices.Thermostats[device].HasLeaf)
			fmt.Printf("Mode:\t\t%v\n", devices.Thermostats[device].HvacMode)
			fmt.Printf("Temperature:\t%v\n", devices.Thermostats[device].AmbientTemperatureF)
			fmt.Printf("Humidity:\t%v\n", devices.Thermostats[device].Humidity)
			fmt.Printf("State:\t\t%v\n\n", devices.Thermostats[device].HvacState)
		}
		data, err := client.NewPoint(devices.Thermostats[device].DeviceID, nil, dataPoints, time)
		if err != nil {
			log.Fatalln("Unable to create NewPoint: ", err)
		}
		batchPoints.AddPoint(data)
		influxdbClient.Write(batchPoints)
		if err != nil {
			log.Fatal("Unable to write data to influxdb: ", err)
		}
	}
}
