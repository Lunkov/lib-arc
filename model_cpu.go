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

var memCPU = make(map[string]CPU)

func CPUCount() int64 {
  return int64(len(memCPU))
}

func CPUAppend(code string, info *CPU) {
  memCPU[code] = *info
}

func GetCPUByCODE(code string) (*CPU) {
  item, ok := memCPU[code]
  if ok {
    return &item
  }
  return nil
}

func GetCPUFactor(cpu1 string, cpu2 string) float32 {
  item1, ok := memCPU[cpu1]
  if !ok {
    return 0
  }
  item2, ok := memCPU[cpu2]
  if !ok {
    return 0
  }
  return item2.AverageCPUMark / item1.AverageCPUMark
}

func LoadCPUsFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false {
      if glog.V(2) {
        glog.Infof("LOG: CPU file: %s", filename)
      }
      var err error
      yamlFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.CPU(%s)  #%v ", filename, err)
      } else {
        count += fileCPUParse(filename, yamlFile)
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

func fileCPUParse(filename string, yamlFile []byte) int {
  var err error
  oTmp := make(map[string]CPU)

  err = yaml.Unmarshal(yamlFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: EntityFile(%s): YAML: %v", filename, err)
    return 0
  }
  
  for key, value := range oTmp {
    value.Name = key
    memCPU[key] = value
  }

  return len(oTmp)
}
