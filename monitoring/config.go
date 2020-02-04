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
