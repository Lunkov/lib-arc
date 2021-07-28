package hardware

import (
  "os"
  "sort"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Motherboard struct {
  Name             string      `db:"name"           json:"name"            yaml:"name"                         unique:"true"`
  Type             string      `db:"type"           json:"type"            yaml:"type"                         unique:"type"`
  SizeU            int         `db:"size_u"         json:"size_u"          yaml:"size_u"                       unique:"size_u"`
  Vendor           string      `db:"vendor"         json:"vendor"          yaml:"vendor"                       unique:"vendor"`
  MemoryMax        int         `db:"memory"         json:"memory"          yaml:"memory"                       unique:"memory"`
  Links            []string    `db:"links"          json:"links"           yaml:"links"                        unique:"links"`
}

type MotherboardSlice []*Motherboard

// Len is part of sort.Interface.
func (d MotherboardSlice) Len() int {
  return len(d)
}

// Swap is part of sort.Interface.
func (d MotherboardSlice) Swap(i, j int) {
  d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d MotherboardSlice) Less(i, j int) bool {
  return d[i].Name < d[j].Name
}

type MotherboardSet struct {
  a      map[string]Motherboard
  index  MotherboardSlice
}

var extMotherboard = ".node"

func NewMotherboards() *MotherboardSet {
  return &MotherboardSet{
                 a: make(map[string]Motherboard),
                 index: make(MotherboardSlice, 0, 0),
               }
}

func (s *MotherboardSet) Count() int64 {
  return int64(len(s.a))
}

func (s *MotherboardSet) Append(info Motherboard) {
  _, ok := s.a[info.Name]
  if ok {
    glog.Warningf("WRN: Motherboard Exists: %s", info.Name)
  }
  s.a[info.Name] = info
  s.index = append(s.index, &info)
  sort.Sort(&s.index)
}

func (s *MotherboardSet) GetByCODE(code string) (*Motherboard) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *MotherboardSet) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == extMotherboard  {
      if glog.V(2) {
        glog.Infof("LOG: Motherboard file: %s", filename)
      }
      var err error
      yamlFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Motherboard(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, yamlFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Motherboards: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *MotherboardSet) fileParse(filename string, yamlFile []byte) int {
  var err error
  oTmp := make(map[string]Motherboard)

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
