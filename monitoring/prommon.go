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
package monitoring

import (
  "log"
  "github.com/prometheus/client_golang/prometheus"
  "time"
  "runtime"
)

var (
  latency = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "zookeeper_monitor_latency",
      Help: "Zookeeper monitor latency info",
    },
    []string{"cluster_name", "host_address", "latency_type", "server_state"})

  outstanding_requests = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "zookeeper_monitor_outstanding_requests",
      Help: "Zookeeper monitor outstanding requests info",
    },
    []string{"cluster_name", "host_address", "server_state"})

  packets = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "zookeeper_monitor_packets",
      Help: "Zookeeper monitor packets info",
    },
    []string{"cluster_name", "host_address", "direction", "server_state"})

  fd = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "zookeeper_monitor_file_descriptor_count",
      Help: "Zookeeper monitor file descriptorcount info",
    },
    []string{"cluster_name", "host_address", "type", "server_state"})

  nodecnt = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "zookeeper_monitor_znode_count",
      Help: "Zookeeper monitor znode count info",
    },
    []string{"cluster_name", "host_address", "server_state"})

  watchcnt = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "zookeeper_monitor_watch_count",
      Help: "Zookeeper monitor watch count info",
    },
    []string{"cluster_name", "host_address", "server_state"})

  follower = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "zookeeper_monitor_follower",
      Help: "Zookeeper monitor follower info",
    },
    []string{"cluster_name", "type"})
)


func RegisterMonitors(registry *prometheus.Registry) {
  registry.MustRegister(latency)
  registry.MustRegister(outstanding_requests)
  registry.MustRegister(packets)
  registry.MustRegister(fd)
  registry.MustRegister(nodecnt)
  registry.MustRegister(watchcnt)
  registry.MustRegister(follower)
}

func MonitorCluster(config *Config,cls ZKCluster) {
  for {
    cmr, err := cls.Monitor()

    log.Print("Cluster metrics collected: " + cls.Name)
    if ( err != nil ) {
      log.Fatal("At least one error at cluster monitoring", err)
    }

    for h, data := range cmr {
      latency.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "latency_type": "zk_avg_latency",
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_avg_latency"].(float64))

      latency.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "latency_type": "zk_min_latency",
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_min_latency"].(float64))

      latency.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "latency_type": "zk_max_latency",
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_max_latency"].(float64))

      outstanding_requests.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_outstanding_requests"].(float64))

      packets.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "direction": "sent",
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_packets_received"].(float64))

      packets.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "direction": "received",
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_packets_received"].(float64))

      fd.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "type": "open",
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_open_file_descriptor_count"].(float64))

      fd.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "type": "max",
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_max_file_descriptor_count"].(float64))

      nodecnt.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_znode_count"].(float64))

      watchcnt.With(prometheus.Labels{"cluster_name": cls.Name,
        "host_address": h,
        "server_state": data["zk_server_state"].(string)}).Set(data["zk_watch_count"].(float64))

      if( data["zk_server_state"].(string) == "leader" ) {
        follower.With(prometheus.Labels{"cluster_name": cls.Name,
          "type": "pending"}).Set(data["zk_pending_syncs"].(float64))

        follower.With(prometheus.Labels{"cluster_name": cls.Name,
          "type": "count"}).Set(data["zk_followers"].(float64))

        follower.With(prometheus.Labels{"cluster_name": cls.Name,
          "type": "synced"}).Set(data["zk_synced_followers"].(float64))
      }
    }

    time.Sleep(time.Duration(config.QueryTime) * time.Second)
  }
}

func MonitorAll(config *Config) {
  runtime.GOMAXPROCS(len(config.Clusters)+2)
  for _, cls := range config.Clusters {
    log.Print("Monitor starting for " + cls.Name)
    go MonitorCluster(config, cls)
  }

}
