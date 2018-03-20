package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

var (
	onewireDevicePath   = "/sys/bus/w1/devices/"
	onewireDeviceList   []string
	listenAddress       = flag.String("web.listen", ":8105", "Address and port to expose metrics")
	metricsPath         = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	onewireTemperatureC = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onewire_temperature_c",
			Help: "Onewire Temperature Sensor Value in Celsius.",
		},
		[]string{
			"device_id",
		},
	)
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
	// Parsing command line arguments
	flag.Parse()
	// Registers temperature gauges
	prometheus.MustRegister(onewireTemperatureC)
}

func main() {
	log.Info("Started")
	// install prometheus http handler
	http.Handle(*metricsPath, prometheus.Handler())
	// TODO: add friendly handling of / path (now returns 404 with a blank page)
	// launch prometheus metrics handler as a goroutine
	go observeOnewireTemperature()
	// starts http listener
	log.WithFields(log.Fields{"httpListen": *listenAddress}).Info("Exporter listening")
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}

func observeOnewireTemperature() {
	// lists onewire devices
	err := createOnewireDeviceList()
	if err != nil {
		log.Fatal("Error getting Onewire device list")
	}
	for {
		for _, deviceID := range onewireDeviceList {
			value, err := readOnewireDevicePayload(deviceID)
			if err != nil {
				log.WithFields(log.Fields{"deviceID": deviceID}).Error("Error reading from device")
			}
			log.WithFields(log.Fields{"deviceID": deviceID, "value": value}).Info("Value read from device")
			onewireTemperatureC.With(prometheus.Labels{"device_id": deviceID}).Set(value)
		}
		time.Sleep(60 * time.Second)
	}
}

func readOnewireDevicePayload(deviceID string) (float64, error) {
	devicePayloadFile := fmt.Sprintf("%s%s/w1_slave", onewireDevicePath, deviceID)
	buffer, err := ioutil.ReadFile(devicePayloadFile)
	if err != nil {
		log.WithFields(log.Fields{"devicePayloadFile": devicePayloadFile}).Error("Error reading Device")
		return 0, err
	}
	re := regexp.MustCompile(`(?s).*YES.*t=([0-9]+)`)
	value, _ := strconv.ParseFloat(re.FindStringSubmatch(string(buffer))[1], 64)
	return value / 1000, nil
}

func createOnewireDeviceList() error {
	devices, err := ioutil.ReadDir(onewireDevicePath)
	if err != nil {
		log.Fatalf("Can't read device directory %v", err)
		return nil
	}
	// searching for onewire attached devices
	for _, device := range devices {
		if strings.Contains(device.Name(), "w1_bus_master1") != true {
			onewireDeviceList = append(onewireDeviceList, device.Name())
			log.Infof("Device found: %s", device.Name())
		}
	}
	return nil
}
