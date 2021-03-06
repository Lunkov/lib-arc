package arc

import (
  "os"
  "sort"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type CPU struct {
  Name             string `db:"name"           json:"name"            yaml:"name"                         unique:"true"`
  
  AverageCPUMark   float32  `db:"average-cpu-mark"  json:"average-cpu-mark" yaml:"average-cpu-mark"`
  Threads          int      `db:"threads"           json:"threads"          yaml:"threads"`
  Cores            int      `db:"cores"             json:"cores"            yaml:"cores"`
}

type CPUSlice []*CPU

// Len is part of sort.Interface.
func (d CPUSlice) Len() int {
  return len(d)
}

// Swap is part of sort.Interface.
func (d CPUSlice) Swap(i, j int) {
  d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d CPUSlice) Less(i, j int) bool {
  return d[i].Name < d[j].Name
}

type CPUSet struct {
  a      map[string]CPU
  index  CPUSlice
}

var extCPU = ".cpu"

func NewCPUs() *CPUSet {
  return &CPUSet{
                 a: make(map[string]CPU),
                 index: make(CPUSlice, 0, 0),
               }
}

func (s *CPUSet) Count() int64 {
  return int64(len(s.a))
}

func (s *CPUSet) Append(info CPU) {
  _, ok := s.a[info.Name]
  if ok {
    glog.Warningf("WRN: CPU Exists: %s", info.Name)
  }
  s.a[info.Name] = info
  s.index = append(s.index, &info)
  sort.Sort(&s.index)
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
    s.Append(value)
  }

  return len(oTmp)
}
