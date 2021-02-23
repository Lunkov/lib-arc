package main

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Property struct {
  CODE         string    `db:"code"           json:"code"            yaml:"code"`
  Name         string    `db:"name"           json:"name"            yaml:"name"`
  Type         string    `db:"type"           json:"type"            yaml:"type"`
}

type Data struct {
  CODE         string     `db:"code"           json:"code"            yaml:"code"`
  Name         string     `db:"name"           json:"name"            yaml:"name"`
  Count        uint64     `db:"count"          json:"count"           yaml:"count"`
  Props      []Property   `db:"properties"     json:"properties"      yaml:"properties"`
}

var memData = make(map[string]Data)

func DataCount() int64 {
  return int64(len(memData))
}

func DataAppend(info *Data) {
  memData[info.CODE] = *info
}

func GetDataByCODE(code string) (*Data) {
  item, ok := memData[code]
  if ok {
    return &item
  }
  return nil
}

func LoadDataFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false {
      if glog.V(2) {
        glog.Infof("LOG: Data file: %s\n", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Data(%s)  #%v ", filename, err)
      } else {
        count += fileDataParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Entities: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func fileDataParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Data

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: DataFile(%s): JSON: %v", filename, err)
    return 0
  }
  DataAppend(&oTmp)
  return 1
}

func RenderDiaClass(code string) string {
  res := "classdiagram"
  
  return res
}
