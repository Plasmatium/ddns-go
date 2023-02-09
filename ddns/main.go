package main

import (
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// set logger
func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

// get record list from env
var recordList []string

func init() {
	recordStr := os.Getenv("RECORD_LIST")
	recordList = strings.Split(recordStr, ",")
	log.Infof("records: %v", recordList)
}

func main() {
	daemonEnvStr := strings.ToLower(os.Getenv("RUN_AS_DAEMON"))
	if daemonEnvStr == "true" {
		log.Info("run in daemon mode")
		runDaemon()
	} else {
		log.Info("run once")
		runOnce()
	}
}

func runOnce() {
	for _, r := range recordList {
		TrySetDNS(r)
	}
}

func runDaemon() {
	// zero time point hit
	runOnce()

	tickChan := time.Tick(time.Minute * 10)
	for {
		<-tickChan
		runOnce()
	}
}
