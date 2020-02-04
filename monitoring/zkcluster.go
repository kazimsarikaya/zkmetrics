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
)

type ZKCluster struct {
    Name string `yaml:"name"`
    Hosts []ZKHost `yaml:hosts`
}

func (cls *ZKCluster) Reset() error {
  for _, h := range cls.Hosts {
    err := h.Reset()
    if (err != nil) {
      log.Print("Cannot reset host metrics: " + h.Address + " ", err)
    }
  }
  return nil
}

func (cls *ZKCluster) Monitor() (map[string]map[string]interface{}, error){
  result := make(map[string]map[string]interface{})
  var err error = nil
  for _, h := range cls.Hosts {
    mr, err := h.Monitor()
    if err != nil {
      log.Print("Can not monitor zk host: " + h.Address, err)
    } else {
      result[h.Address] = mr
    }
  }

  return result, err
}
