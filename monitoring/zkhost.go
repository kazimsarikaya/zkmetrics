package monitoring

import (
  "net"
  "io/ioutil"
  "bufio"
  "time"
  "bytes"
  "strings"
  "strconv"

  "github.com/pkg/errors"
)

type ZKHost struct {
  Address string `yaml:"address"`
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
    if parts[0] == "zk_version" || parts[0] == "zk_server_state" {
      result[parts[0]]=parts[1]
    } else {
      t, _ := strconv.Atoi(parts[1])
      result[parts[0]] = t
    }
  }

  return result, nil
}
