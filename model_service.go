package main

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Service struct {
  CODE         string    `db:"code"           json:"code,omitempty"            yaml:"code"`
  Type         string    `db:"type"           json:"type,omitempty"            yaml:"type"`
  Name         string    `db:"name"           json:"name,omitempty"            yaml:"name"`
  SystemCode   string    `db:"system_code"    json:"system_code,omitempty"     yaml:"system_code"`
  Disabled     bool      `db:"disabled"       json:"disabled,omitempty"        yaml:"disabled"`
}

var memService = make(map[string]Service)

func ServiceCount() int64 {
  return int64(len(memService))
}

func ServiceAppend(info *Service) {
  memService[info.CODE] = *info
}

func GetServiceByCODE(code string) (*Service) {
  item, ok := memService[code]
  if ok {
    return &item
  }
  return nil
}

func LoadServicesFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == ".service"  {
      if glog.V(2) {
        glog.Infof("LOG: Service file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Service(%s)  #%v ", filename, err)
      } else {
        count += fileServiceParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Services: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func fileServiceParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Service

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: ServiceFile(%s): JSON: %v", filename, err)
  }
  if !oTmp.Disabled {
    ServiceAppend(&oTmp)
    return 1
  }

  return 0
}


