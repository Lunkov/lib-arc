package arc

import (
  "strings"
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Role struct {
  CODE         string                   `db:"code"           json:"code,omitempty"            yaml:"code"`
  Type         string                   `db:"type"           json:"type,omitempty"            yaml:"type"`
  Name         string                   `db:"name"           json:"name,omitempty"            yaml:"name"`
  Description  string                   `db:"description"    json:"description,omitempty"     yaml:"description"`
}

type Roles struct {
  a map[string]Role
}

var extRole = ".role"

func NewRoles() *Roles {
  return &Roles{
                 a: make(map[string]Role),
               }
}

func (s *Roles) Count() int64 {
  return int64(len(s.a))
}

func (s *Roles) Append(info *Role) {
  info.CODE = strings.ToLower(info.CODE)
  s.a[info.CODE] = *info
}

func (s *Roles) GetByCODE(code string) (*Role) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Roles) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == extRole  {
      if glog.V(2) {
        glog.Infof("LOG: Role file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Role(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Roles: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *Roles) fileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Role

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: RoleFile(%s): JSON: %v", filename, err)
  }
  s.Append(&oTmp)
  return 1
}
