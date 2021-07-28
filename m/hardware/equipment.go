package hardware

import (
  "os"
  "sort"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Equipment struct {
  Name             string      `db:"name"           json:"name"            yaml:"name"                         unique:"true"`
  Type             string      `db:"type"           json:"type"            yaml:"type"                         unique:"type"`
  SizeU            int         `db:"size_u"         json:"size_u"          yaml:"size_u"                       unique:"size_u"`
  Vendor           string      `db:"vendor"         json:"vendor"          yaml:"vendor"                       unique:"vendor"`
  Include          []string    `db:"include"        json:"include"         yaml:"include"                      unique:"include"`
  Links            []string    `db:"links"          json:"links"           yaml:"links"                        unique:"links"`
}

type EquipmentSlice []*Equipment

// Len is part of sort.Interface.
func (d EquipmentSlice) Len() int {
  return len(d)
}

// Swap is part of sort.Interface.
func (d EquipmentSlice) Swap(i, j int) {
  d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d EquipmentSlice) Less(i, j int) bool {
  return d[i].Name < d[j].Name
}

type EquipmentSet struct {
  a      map[string]Equipment
  index  EquipmentSlice
}

var extEquipment = ".node"

func NewEquipments() *EquipmentSet {
  return &EquipmentSet{
                 a: make(map[string]Equipment),
                 index: make(EquipmentSlice, 0, 0),
               }
}

func (s *EquipmentSet) Count() int64 {
  return int64(len(s.a))
}

func (s *EquipmentSet) Append(info Equipment) {
  _, ok := s.a[info.Name]
  if ok {
    glog.Warningf("WRN: Equipment Exists: %s", info.Name)
  }
  s.a[info.Name] = info
  s.index = append(s.index, &info)
  sort.Sort(&s.index)
}

func (s *EquipmentSet) GetByCODE(code string) (*Equipment) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *EquipmentSet) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == extEquipment  {
      if glog.V(2) {
        glog.Infof("LOG: Equipment file: %s", filename)
      }
      var err error
      yamlFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Equipment(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, yamlFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Equipments: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *EquipmentSet) fileParse(filename string, yamlFile []byte) int {
  var err error
  oTmp := make(map[string]Equipment)

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
