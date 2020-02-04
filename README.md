# Multiple Zookeeper Prometheus Metric Exporter

This program monitors multiple **zk cluster**s or single **zk**.

It uses a configuration file with format:

```
version: "v0.1" # arbitrary version
name: "zkmonitor1" # arbitrary name
queryTime: 5 # how many second each gather
resetAfter: 100 # reset zk metrics on zk side
clusters: # one or more cluster
- name: "zkcluster1" # cluster name
  hosts: # one or more host
  - address: "zkhost:2181"
```
