package main

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "encoding/json"
)

type SwaggerServers struct {
  URL         string    `db:"url"           json:"url"            yaml:"url"`
}

type SwaggerInfo struct {
  Title         string    `db:"title"           json:"title"            yaml:"title"`
  Version       string    `db:"version"         json:"version"          yaml:"version"`
  Description   string    `db:"description"     json:"description"      yaml:"description"`
}

type SwaggerMethod struct {
  OperationId   string    `db:"operationId"           json:"operationId"            yaml:"operationId"`
  Summary       string    `db:"summary"               json:"summary"                yaml:"summary"`
  Description   string    `db:"description"           json:"description"            yaml:"description"`
}

type SwaggerProperty struct {
  Type         string    `db:"type"           json:"type"            yaml:"type"`
  Format       string    `db:"format"         json:"format"          yaml:"format"`
  Description   string   `db:"description"    json:"description"     yaml:"description"`
}

type SwaggerMethods map[string]SwaggerMethod

type Swagger struct {
  Name         string            `db:"swagger"     json:"swagger"         yaml:"swagger"`
  Host         string            `db:"host"        json:"host"            yaml:"host"`
  Info         SwaggerInfo       `db:"info"        json:"info"            yaml:"info"`
  Servers      []SwaggerServers  `db:"servers"     json:"servers"         yaml:"servers"`
  Schemes      []string          `db:"schemes"     json:"schemes"         yaml:"schemes"`
  Consumes     []string                        `db:"consumes"       json:"consumes"        yaml:"consumes"`
  Produces     []string                        `db:"produces"       json:"produces"        yaml:"produces"`
  Paths        map[string]SwaggerMethods       `db:"paths"       json:"paths"        yaml:"paths"`
}

type SwaggerSet struct {
  a map[string]map[string]Swagger // map[Title]map[Version]
}

func NewSwaggerSet() *SwaggerSet {
  return &SwaggerSet{
                 a: make(map[string]map[string]Swagger),
               }
}

func (s *SwaggerSet) Count() int64 {
  return int64(len(s.a))
}

func (s *SwaggerSet) Append(info *Swagger) {
  if _, ok := s.a[info.Info.Title]; !ok {
    s.a[info.Info.Title] = make(map[string]Swagger)
  }

  s.a[info.Info.Title][info.Info.Version] = (*info)
}

func (s *SwaggerSet) GetList() []SwaggerInfo {
  res := make([]SwaggerInfo, 0)
  for _, versions := range s.a {
    for _, value := range versions {
      res = append(res, value.Info)
    }
  }
  return res
}

func (s *SwaggerSet) GetByTitle(code string, version string) (*Swagger) {
  item, ok := s.a[code][version]
  if ok {
    return &item
  }
  itemT, ok0 := s.a[code]
  if ok0 {
    for _, item0 := range itemT {
      return &item0
    }
  }
  return nil
}

func (s *SwaggerSet) AppendToService(service *Service, code string, version string) {
  
}

func (s *SwaggerSet) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false {
      if glog.V(2) {
        glog.Infof("LOG: Swagger file: %s\n", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Swagger(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Items: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *SwaggerSet) fileParse(filename string, jsonFile []byte) int {
  var err error
  var swgTmp Swagger

  err = json.Unmarshal(jsonFile, &swgTmp)
  if err != nil {
    glog.Errorf("ERR: swaggerFile(%s): JSON: %v", filename, err)
  }
  s.Append(&swgTmp)

  return 1
}
