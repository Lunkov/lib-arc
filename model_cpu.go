package main

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type CPU struct {
  Name             string `db:"name"           json:"name"            yaml:"name"`
  
  AverageCPUMark   float32  `db:"average-cpu-mark"  json:"average-cpu-mark" yaml:"average-cpu-mark"`
  Threads          int      `db:"threads"           json:"threads"          yaml:"threads"`
  Cores            int      `db:"cores"             json:"cores"            yaml:"cores"`
}

type CPUSet struct {
  a map[string]CPU
}

var extCPU = ".cpu"

func NewCPUs() *CPUSet {
  return &CPUSet{
                 a: make(map[string]CPU),
               }
}

func (s *CPUSet) Count() int64 {
  return int64(len(s.a))
}

func (s *CPUSet) Append(code string, info *CPU) {
  s.a[code] = *info
}

func (s *CPUSet) GetByCODE(code string) (*CPU) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *CPUSet) GetFactor(cpu1 string, cpu2 string) float32 {
  item1, ok := s.a[cpu1]
  if !ok {
    return 0
  }
  item2, ok := s.a[cpu2]
  if !ok {
    return 0
  }
  return item2.AverageCPUMark / item1.AverageCPUMark
}

func (s *CPUSet) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == extCPU  {
      if glog.V(2) {
        glog.Infof("LOG: CPU file: %s", filename)
      }
      var err error
      yamlFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.CPU(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, yamlFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, CPUs: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *CPUSet) fileParse(filename string, yamlFile []byte) int {
  var err error
  oTmp := make(map[string]CPU)

  err = yaml.Unmarshal(yamlFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: EntityFile(%s): YAML: %v", filename, err)
    return 0
  }
  
  for key, value := range oTmp {
    value.Name = key
    s.a[key] = value
  }

  return len(oTmp)
}
