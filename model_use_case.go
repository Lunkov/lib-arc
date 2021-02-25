package arc

import (
  "os"
  "sort"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type UseCase struct {
  CODE           string    `db:"code"           json:"code,omitempty"            yaml:"code"`
  Type           string    `db:"type"           json:"type,omitempty"            yaml:"type"`
  Path           string    `db:"path"           json:"-"                         yaml:"-"`
  Date           string    `db:"date"           json:"date,omitempty"            yaml:"date"`
  Name           string    `db:"name"           json:"name,omitempty"            yaml:"name"`
  Disabled       bool      `db:"disabled"       json:"disabled,omitempty"        yaml:"disabled"`
  URL            string    `db:"url"            json:"url,omitempty"             yaml:"url"`
  Services     []string    `db:"services"       json:"services,omitempty"        yaml:"services"`
  Systems      []string    `db:"systems"        json:"systems,omitempty"         yaml:"systems"`
  Sequences    []string    `db:"sequences"      json:"sequences,omitempty"       yaml:"sequences"`
  Tasks        []string    `db:"tasks"          json:"tasks,omitempty"           yaml:"tasks"`
  ReadMe         string    `db:"readme"         json:"readme,omitempty"          yaml:"-"`
}

var extUseCase = ".case"

type useCaseSlice []*UseCase

// Len is part of sort.Interface.
func (d useCaseSlice) Len() int {
  return len(d)
}

// Swap is part of sort.Interface.
func (d useCaseSlice) Swap(i, j int) {
  glog.Warningf("WRN: Use Case Swap: i=%d j=%d", i, j)
  d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d useCaseSlice) Less(i, j int) bool {
  glog.Warningf("WRN: Use Case Less: i=%d j=%d", i, j)
  return d[i].Date < d[j].Date
}

type UseCaseSet struct {
  a      map[string]UseCase
  index  useCaseSlice
}

func NewUseCaseSet() *UseCaseSet {
  return &UseCaseSet{
                 a: make(map[string]UseCase),
                 index: make(useCaseSlice, 0, 0),
               }
}

func (s *UseCaseSet) Count() int64 {
  return int64(len(s.a))
}

func (s *UseCaseSet) Append(info UseCase) {
  _, ok := s.a[info.CODE]
  if ok {
    glog.Warningf("WRN: Case Exists: %s", info.CODE)
  }
  s.a[info.CODE] = info
  s.index = append(s.index, &info)
  sort.Sort(&s.index)
}

func (s *UseCaseSet) GetByCODE(code string) (*UseCase) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *UseCaseSet) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == extUseCase {
      if glog.V(2) {
        glog.Infof("LOG: Case file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Case(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Find Cases: %d", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }
  
  /* TODO
   * for k, item := range s.index {
    glog.Infof("LOG: Scan Path For Case(%d) '%s': %s", k, item.CODE, item.Path)
    LoadSequencesFromFiles(item.Path)
  }*/
  return count
}

func (s *UseCaseSet) fileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp UseCase

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: CaseFile(%s): JSON: %v", filename, err)
  }
  if !oTmp.Disabled {
    oTmp.Path = filepath.Dir(filename)
    s.Append(oTmp)
    return 1
  }

  return 0
}
