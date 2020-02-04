package monitoring

import (
  "github.com/pkg/errors"
)

type ZKCluster struct {
    Name string `yaml:"name"`
    Hosts []ZKHost `yaml:hosts`
}


func (cls *ZKCluster) Monitor() (map[string]map[string]interface{}, error){
  result := make(map[string]map[string]interface{})
  for _, h := range cls.Hosts {
    mr, err := h.Monitor()
    if err != nil {
      return nil, errors.Wrapf(err, "Can not monitor zk host: " + h.Address)
    }
    result[h.Address] = mr
  }

  return result, nil
}
