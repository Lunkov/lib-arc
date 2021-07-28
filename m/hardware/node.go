package hardware

import (
  "os"
  "sort"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Node struct {
  Name             string      `db:"name"           json:"name"            yaml:"name"                         unique:"true"`
  Type             string      `db:"type"           json:"type"            yaml:"type"                         unique:"type"`
  SizeU            int         `db:"size_u"         json:"size_u"          yaml:"size_u"                       unique:"size_u"`
  Vendor           string      `db:"vendor"         json:"vendor"          yaml:"vendor"                       unique:"vendor"`
  Memory           int         `db:"memory"         json:"memory"          yaml:"memory"                       unique:"memory"`
  Links            []string    `db:"links"          json:"links"           yaml:"links"                        unique:"links"`
  TotalCost        float32     `db:"total_cost"     json:"total_cost"      yaml:"total_cost"                        unique:"links"`
}

type NodeSlice []*Node

// Len is part of sort.Interface.
func (d NodeSlice) Len() int {
  return len(d)
}

// Swap is part of sort.Interface.
func (d NodeSlice) Swap(i, j int) {
  d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d NodeSlice) Less(i, j int) bool {
  return d[i].Name < d[j].Name
}

type NodeSet struct {
  a      map[string]Node
  index  NodeSlice
}

var extNode = ".node"

func NewNodes() *NodeSet {
  return &NodeSet{
                 a: make(map[string]Node),
                 index: make(NodeSlice, 0, 0),
               }
}

func (s *NodeSet) Count() int64 {
  return int64(len(s.a))
}

func (s *NodeSet) Append(info Node) {
  _, ok := s.a[info.Name]
  if ok {
    glog.Warningf("WRN: Node Exists: %s", info.Name)
  }
  s.a[info.Name] = info
  s.index = append(s.index, &info)
  sort.Sort(&s.index)
}

func (s *NodeSet) GetByCODE(code string) (*Node) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *NodeSet) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == extNode  {
      if glog.V(2) {
        glog.Infof("LOG: Node file: %s", filename)
      }
      var err error
      yamlFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Node(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, yamlFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Nodes: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *NodeSet) fileParse(filename string, yamlFile []byte) int {
  var err error
  oTmp := make(map[string]Node)

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
