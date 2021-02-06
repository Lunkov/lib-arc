package main

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type System struct {
  CODE         string    `db:"code"           json:"code,omitempty"            yaml:"code"`
  Type         string    `db:"type"           json:"type,omitempty"            yaml:"type"`
  Name         string    `db:"name"           json:"name,omitempty"            yaml:"name"`
  Disabled     bool      `db:"disabled"       json:"disabled,omitempty"        yaml:"disabled"`
}

var memSystem = make(map[string]System)

func SystemCount() int64 {
  return int64(len(memSystem))
}

func SystemAppend(info *System) {
  memSystem[info.CODE] = *info
}

func GetSystemByCODE(code string) (*System) {
  item, ok := memSystem[code]
  if ok {
    return &item
  }
  return nil
}

func LoadSystemsFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == ".system"  {
      if glog.V(2) {
        glog.Infof("LOG: System file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.System(%s)  #%v ", filename, err)
      } else {
        count += fileSystemParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Systems: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func fileSystemParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp System

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: SystemFile(%s): JSON: %v", filename, err)
  }
  if !oTmp.Disabled {
    SystemAppend(&oTmp)
    return 1
  }

  return 0
}


