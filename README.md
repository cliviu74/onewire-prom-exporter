# Prometheus OneWire sensors exporter

This is an exporter that exposes information gathered from OneWire
sensors in Prometheus friendly format.

## Prerequisites

Before building this exporter you need a working golang development environment. A good starting point is the 
[Golang Getting Started](https://golang.org/doc/install) document.
## Building and Installation

The current build.sh script will build the binary for ARM 5 architecture (good for Raspberry Pi)
```
go get github.com/cliviu74/onewire-prom-exporter
cd $GOPATH/src/github.com/cliviu74/onewire-prom-exporter
./build.sh
```

## Usage

```
./onewire-prom-exporter
```

Visit http://localhost:8105/metrics to get metrics from exporter. The exporter will list all sensors
and display the metrics as gauges, labeled with device_id (onwire address)

```
~# curl -s http://localhost:8105/metrics | grep onewire
# HELP onewire_temperature_c Onewire Temperature Sensor Value in Celsius.
# TYPE onewire_temperature_c gauge
onewire_temperature_c{device_id="28-xxxxxxxxxxxx",hostname="pienv"} 29.125
onewire_temperature_c{device_id="28-xxxxxxxxxxxx",hostname="pienv"} 31.812

```

## Configuration

You can also specify a the web port and metrics path upon launching the exporter

```
Usage of ./onewire-prom-exporter:
  -web.listen-address string
    	Address and port to expose metrics (default ":8105")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
```

```
./onewire-prom-exporter -web.listen-address=0.0.0.0:8105 -web.telemetry-path=/metrics 
```

## Runing the exporter with systemd

see examples: [here](examples/systemd/README.md)
## Prometheus Configuration

This is a simple target configuration for prometheus. The ip in the targets array
needs to be replaced with the IP of the node running the exporter.

Example config:
```YAML
scrape_configs:
  - job_name: 'onewire_exporter'
    scrape_interval: 60s
    metrics_path: /metrics
    scheme: http  
    static_configs:
      - targets: ['192.168.1.123:8105']
```
## Contributing

Feel free to report bugs, contribute, discuss changes, fork this project, create pull requests.
If you wish to become a contributor, you are welcome to contact me. 
## License
Apache License 2.0, see [LICENSE](LICENSE).
