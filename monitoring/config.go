package monitoring

import (
  "io/ioutil"
  "path/filepath"

  "gopkg.in/yaml.v2"

)

type Config struct {
  Version string `yaml:"version"`
  Name string `yaml:"name"`
  QueryTime int `yaml:"queryTime"`
  Clusters []ZKCluster `yaml:"clusters,omitempty"`
}


func Parse(file string) (*Config, error) {
  filefpath, _ := filepath.Abs(file)
  filecontent, err := ioutil.ReadFile(filefpath)
  if err != nil {
    return nil, err
  }
  config := new(Config)
  err = yaml.Unmarshal(filecontent,config)
  return config, err
}

func (cfg *Config) Print() string {
  result := "version -> " + cfg.Version + ", " + "config name -> " +
    cfg.Name + ", "
  for _, c := range cfg.Clusters {
    result += "(cluster name -> " + c.Name + ", hosts -> "
    for _, h := range c.Hosts {
      result += h.Address + ","
    }
    result += "), "
  }

  return result
}
