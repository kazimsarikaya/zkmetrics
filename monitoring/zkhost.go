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
  "net"
  "io/ioutil"
  "bufio"
  "time"
  "bytes"
  "strings"
  "strconv"
  "log"

  "github.com/pkg/errors"
)

type ZKHost struct {
  Address string `yaml:"address"`
}

func (zk *ZKHost) Reset() error {
  timeout, _ := time.ParseDuration("5s")

  conn, err := net.DialTimeout("tcp", zk.Address, timeout)
  if err != nil {
    return err
  }
  defer conn.Close()

  err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}

  _, err = conn.Write([]byte("srst"))
  if err != nil {
    return errors.Wrapf(err, "cannot send mntr command")
  }

  _, err = conn.Write([]byte("crst"))
  if err != nil {
    return errors.Wrapf(err, "cannot send mntr command")
  }

  return nil
}

func (zk *ZKHost) Monitor() (map[string]interface{}, error) {
  timeout, _ := time.ParseDuration("5s")

  conn, err := net.DialTimeout("tcp", zk.Address, timeout)
  if err != nil {
    return nil, err
  }
  defer conn.Close()

  err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}

  _, err = conn.Write([]byte("mntr"))
  if err != nil {
    return nil, errors.Wrapf(err, "cannot send mntr command")
  }

  data, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, errors.Wrap(err, "mntr command read failed")
	}

  result := make(map[string]interface{})

  s := bufio.NewScanner(bytes.NewReader(data))
  for s.Scan() {
    line := s.Text()
    parts := strings.SplitN(line, "\t", 2)
    if len(parts) != 2 {
      log.Print("zk metric line: bad format: " + line)
    } else {
      if parts[0] == "zk_version" || parts[0] == "zk_server_state" {
        result[parts[0]]=parts[1]
      } else {
        t, _ := strconv.Atoi(parts[1])
        result[parts[0]] = float64(t)
      }
    }
  }

  return result, nil
}
