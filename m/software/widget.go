package software

import (
  "os"
  "io/ioutil"
  "strings"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Widget struct {
  CODE         string    `db:"code"           json:"code,omitempty"            yaml:"code"`
  Type         string    `db:"type"           json:"type,omitempty"            yaml:"type"`
  Name         string    `db:"name"           json:"name,omitempty"            yaml:"name"`
  Disabled     bool      `db:"disabled"       json:"disabled,omitempty"        yaml:"disabled"`
}

type Widgets struct {
  a map[string]Widget
}

func NewWidgets() *Widgets {
  return &Widgets{
                 a: make(map[string]Widget),
               }
}

func (s *Widgets) FileExtension() string {
  return ".widget"
}

func (s *Widgets) Count() int64 {
  return int64(len(s.a))
}

func (s *Widgets) Append(info *Widget) {
  info.CODE = strings.ToLower(info.CODE)
  s.a[info.CODE] = *info
}

func (s *Widgets) GetByCODE(code string) (*Widget) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Widgets) GetList() []Widget {
  res := make([]Widget, 0)
  for _, item := range s.a {
    res = append(res, item)
  }
  return res
}

func (s *Widgets) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == s.FileExtension()  {
      if glog.V(2) {
        glog.Infof("LOG: Service file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Service(%s)  #%v ", filename, err)
      } else {
        count += s.FileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Services: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *Widgets) FileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Widget

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: ServiceFile(%s): JSON: %v", filename, err)
  }
  if !oTmp.Disabled {
    s.Append(&oTmp)
    return 1
  }

  return 0
}
