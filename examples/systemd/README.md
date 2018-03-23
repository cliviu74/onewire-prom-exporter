# Onwire Exporter Systemd Unit

This is the unit file that runs the exporter as a service with systemd

The unit file in this directory is to be put into /lib/systemd/system/onewire-exporter.service. 
Service is running as nobody user. Extra configuration for the daemon can be added to  /etc/sysconfig/node_exporter.
Example
```
# /etc/sysconfig/onewire-exporter
OPTIONS="-web.listen-address=:8105 -web.telemetry-path=/metrics"
```

Once unit file is installed at /lib/systemd/system/onewire-exporter.service, service can be enabled and started.

```
systemctl enable onewire-exporter
systemctl start onewire-exporter
```