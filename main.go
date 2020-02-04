/*
ZK Metrics
Copyright (C) 2020  Kazım SARIKAYA <kazimsarikaya@sanaldiyar.com>

This file is part of ZK Metrics.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
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

	monitoring.RegisterMonitors(r)

  handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	monitoring.Monitor(config)

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
