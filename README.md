# go-kafka-lags-prom-exporter

This simple tool exposes kafka lag values for all the topics in your cluster.

I have a [version of this](https://github.com/drio/ksak/blob/main/exporter.go) in my ksak 
tool that uses the segment.io kafka library. But the code is bad and there is a bug when 
[fetching the offets](https://github.com/drio/ksak/blob/main/lag-logic.go#L44) in the 
library (probably cause by me?) where we start and do not close goroutines. Eventually we
collapse the machine since we consume all the cpu bandwidth. Anyway, if you find the bug 
please let me know.

In the meantime, this tool spawns a goroutine which in turn runs a script that uses the 
kafka tools to get the lag metrics. Then we process it and expose the values via a prometheus
metrics:


```sh
# console 1
> go build main.go && ./main -port "9898" -seconds 5 -cmd "./get_lag.sh"
2022/08/03 10:43:24 Listening on port 9898
2022/08/03 10:43:24 Starting go-routine
2022/08/03 10:43:29 Invalid lag in line:0; -; skipping
2022/08/03 10:43:29 Sleeping goroutine for 5 seconds

# console 2
> curl -s http://localhost:9898/metrics | grep kafka_lag
# HELP kafka_lag_exporter lag metrics on kafka topics
# TYPE kafka_lag_exporter gauge
kafka_lag_exporter{group="drio-topic1-gid-1",partition="0",topic="drio-topic1"} 0
kafka_lag_exporter{group="drio-topic2-gid-2",partition="0",topic="drio-topic2"} 2
kafka_lag_exporter{group="drio-topic3-gid-3",partition="0",topic="drio-topic3"} 0
kafka_lag_exporter{group="drio-topic4-gid-4",partition="0",topic="drio-topic4"} 1
kafka_lag_exporter{group="drio-topic5-gid-5",partition="0",topic="drio-topic5"} 1
```

# Running via docker

You probably want to run this via docker. We provide a Dockerfile for that purpose.
Also, see the makefile targets to help you run and test things:

```sh
> make -n docker/run
GOOS=linux GOARCH=amd64 go build -o go-kafka-lags-prom-exporter.linux.amd64
GOOS=darwin GOARCH=arm64 go build -o go-kafka-lags-prom-exporter.darwin.arm64
docker build -t drio-go-kafka-lags-prom-exporter .
docker run -d --name kafka-lags-prom-exporter -p 9898:9898 drio-go-kafka-lags-prom-exporter
```

The cmd runs the ./run.sh script. That script hardcodes a few things you may want to 
modify for your needs.
