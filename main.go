package main

import (
	"log"
	"net/http"
	"time"

	"github.com/kmjayadeep/restic-monitoring/internal/config"
	"github.com/kmjayadeep/restic-monitoring/internal/stats"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LastSnapshot struct {
	Time     time.Time `json:"time"`
	ShortID  string    `json:"short_id"`
	HostName string    `json:"hostname"`
}

func main() {
	c, err := config.ParseConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	go stats.Run(c)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":18090", nil))
}
