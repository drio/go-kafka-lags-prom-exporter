TAG=drio-go-kafka-lags-prom-exporter
PORT=9898

ps:
	@docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Networks}}\t{{.State}}\t{{.CreatedAt}}"

docker/stop:
	docker stop kafka-lags-prom-exporter

docker/bash:
	docker run -it --entrypoint bash $(TAG)

docker/run:	docker/build
	docker run -d --name kafka-lags-prom-exporter -p $(PORT):$(PORT) $(TAG) 

docker/build: build
	docker build -t $(TAG) .

build: go-kafka-lags-prom-exporter.linux.amd64 \
	go-kafka-lags-prom-exporter.darwin.arm64

go-kafka-lags-prom-exporter.linux.amd64:
	GOOS=linux GOARCH=amd64 go build -o go-kafka-lags-prom-exporter.linux.amd64

go-kafka-lags-prom-exporter.darwin.arm64:
	GOOS=darwin GOARCH=arm64 go build -o go-kafka-lags-prom-exporter.darwin.arm64

clean:
	rm -f go-kafka-lags-prom-exporter*
