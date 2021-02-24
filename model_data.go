package main

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
  "sort"
)

// TypeOfData.
type TypeOfData uint64

const (
  TABLE TypeOfData = iota
  VIEW
  PARAMETER
)

var typesOfData = map[string]TypeOfData{
  "table":        TABLE,
  "view":         VIEW,
  "parameter":    PARAMETER,
}

type Property struct {
  CODE         string     `db:"code"           json:"code"            yaml:"code"`
  Name         string     `db:"name"           json:"name"            yaml:"name"`
  GlobalName   string     `db:"global_name"    json:"global_name"     yaml:"global_name"`
  Type         string     `db:"type"           json:"type"            yaml:"type"`
  Order        int        `db:"order"          json:"order"           yaml:"order"`
}

type DataSet struct {
  CODE         string     `db:"code"           json:"code"            yaml:"code"`
  Type         string     `db:"type"           json:"type"            yaml:"type"`
  Name         string     `db:"name"           json:"name"            yaml:"name"`
  GlobalName   string     `db:"global_name"    json:"global_name"     yaml:"global_name"`
  Count        uint64     `db:"count"          json:"count"           yaml:"count"`
  Props      []Property   `db:"properties"     json:"properties"      yaml:"properties"`
}

type DataSets struct {
  a map[string]DataSet
}

func NewDataSets() *DataSets {
  return &DataSets{
                 a: make(map[string]DataSet),
               }
}

func (s *DataSets) Count() int64 {
  return int64(len(s.a))
}

func (s *DataSets) Append(info *DataSet) {
  // Sorting by Order
  sort.Slice(info.Props, func(i, j int) bool { return info.Props[i].Order < info.Props[j].Order })
  s.a[info.CODE] = *info
}

func (s *DataSets) GetByCODE(code string) (*DataSet) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *DataSets) LoadFromFiles(scanPath string) int {
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
        count += s.fileParse(filename, jsonFile)
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

func (s *DataSets) fileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp DataSet

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: DataFile(%s): JSON: %v", filename, err)
    return 0
  }
  s.Append(&oTmp)
  return 1
}

func RenderDiaClass(code string) string {
  res := "classdiagram"
  
  return res
}
