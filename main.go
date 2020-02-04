package main

import (
	"flag"
	"log"
	"net/http"


  "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

  "github.com/kazimsarikaya/zkmetrics/monitoring"
)

var addr = flag.String("listen-address", ":9100", "The address to listen on for HTTP requests.")
var configfile = flag.String("config", "config.yaml", "The configuration file")

func main() {
	flag.Parse()

  log.Print("Parsing config")

  config, err := monitoring.Parse(*configfile)
  if err != nil {
    log.Fatal("Can not parse config: ", err)
  }

  log.Print("Config parsed: " + config.Print())


  r := prometheus.NewRegistry()
  handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
