package monitoring

import (
  "log"
)

type ZKCluster struct {
    Name string `yaml:"name"`
    Hosts []ZKHost `yaml:hosts`
}


func (cls *ZKCluster) Monitor() (map[string]map[string]interface{}, error){
  result := make(map[string]map[string]interface{})
  var err error = nil
  for _, h := range cls.Hosts {
    mr, err := h.Monitor()
    if err != nil {
      log.Fatal("Can not monitor zk host: " + h.Address, err)
    } else {
      result[h.Address] = mr
    }
  }

  return result, err
}
