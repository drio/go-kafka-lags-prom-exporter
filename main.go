package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const version = "0.0.2"

var (
  flgHelp bool
  flgCmd string
  flgPort string
  flgSeconds int
)

// csv lines will be
// GROUP  TOPIC    PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG 
const rows_per_line = 6
type csvEntry struct {
	group string
	topic string
	partition int
	current string
	end string
  lag int
}

func parseCmdLineFlags() {
  flag.BoolVar(&flgHelp, "help", false, "if true, show help")
  flag.StringVar(&flgCmd, "cmd", "", "script to run to get the lag information")
  flag.StringVar(&flgPort, "port", "", "port to listen to")
  flag.IntVar(&flgSeconds, "seconds", 30, "how often to run the lag script")
  flag.Parse()
}

func readCsv(input string) ([]csvEntry, error) {
	reader := csv.NewReader(strings.NewReader(input))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	listEntries := []csvEntry{}
	for idx, entry := range records {
    if len(entry) != rows_per_line {
      log.Printf("Invalid number of columns in csv line:%d", idx)
      continue
    }

    i_partition, err := strconv.Atoi(entry[2])
    if err != nil {
      log.Printf("Invalid partition in line:%d; %s; skipping", idx, entry[2])
      continue
    }
    i_lag, err := strconv.Atoi(entry[5])
    if err != nil {
      log.Printf("Invalid lag in line:%d; %s; skipping", idx, entry[5])
      continue
    }

		listEntries = append(listEntries, csvEntry{
			group: strings.TrimSpace(entry[0]),
      topic: strings.TrimSpace(entry[1]),
      partition: i_partition,
      current: strings.TrimSpace(entry[3]),
      end: strings.TrimSpace(entry[4]),
      lag: i_lag,
		})
	}

	return listEntries, nil
}

func updateGauge(gauge *prometheus.GaugeVec, csvEntries []csvEntry) {
	for _, ce := range csvEntries {
			gauge.With(prometheus.Labels{
				"topic": ce.topic,
				"group": ce.group,
				"partition": fmt.Sprintf("%d",ce.partition),
			}).Set(float64(ce.lag))
	}
}

func registerGauge() *prometheus.GaugeVec {
	gaugeLag := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kafka_lag_exporter",
			Help: "lag metrics on kafka topics",
		},
		[]string{
			"topic",
			"group",
			"partition",
		},
	)
	prometheus.MustRegister(gaugeLag)
	return gaugeLag
}

func runCmd(cmd_string string) (string, error) {
  var stdout, stderr bytes.Buffer

	cmd := exec.Command(cmd_string)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
    return "", err
	}

	err = cmd.Wait()
	if err != nil {
    return "", err
	}
	return string(stdout.Bytes()), nil
}

func main() {
  parseCmdLineFlags()
  if flgHelp {
    flag.Usage()
    os.Exit(0)
  }

  if flgPort == "" {
    flag.Usage()
    os.Exit(0)
  }

  if flgCmd == "" {
    flag.Usage()
    os.Exit(0)
  }

  gaugeLag := registerGauge()

	go func() {
		log.Printf("Starting go-routine")
		for {
      cmdOutput, err := runCmd(flgCmd)
      if err != nil {
        log.Printf("Error running cmd: %s", err)
      } else {
        csvEntries, err := readCsv(cmdOutput)
        if err != nil {
          log.Printf("Error processing cmd output: %s", err)
        } else {
          updateGauge(gaugeLag, csvEntries)
        }
      }
			log.Printf("Sleeping goroutine for %d seconds", flgSeconds)
			time.Sleep(time.Duration(flgSeconds) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Listening on port %s", flgPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", flgPort), nil))
}
