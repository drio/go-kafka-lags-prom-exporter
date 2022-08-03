#!/bin/sh

b="go-kafka-lags-prom-exporter.linux.amd64"
port=9898
seconds=10
./$b -port $port -seconds $seconds -cmd "./get_lag.sh"
