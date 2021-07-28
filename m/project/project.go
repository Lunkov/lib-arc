package project

import (
  "strings"
  "os"
  "time"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Estimate struct {
  CODE         string        `db:"code"           json:"code,omitempty"            yaml:"code"                unique:"true"`

  ProjectCODE  string        `db:"project_code"   json:"project_code,omitempty"    yaml:"project_code"        unique:"true"`
  VariantCODE  string        `db:"variant_code"   json:"variant_code,omitempty"    yaml:"variant_code"        unique:"true"`


  Name         string        `db:"name"           json:"name,omitempty"            yaml:"name"                unique:"true"`
  Description  string        `db:"description"    json:"description,omitempty"     yaml:"description"`
  
  
}

type Milestone struct {
  CODE         string        `db:"code"           json:"code,omitempty"            yaml:"code"                unique:"true"`
  Name         string        `db:"name"           json:"name,omitempty"            yaml:"name"                unique:"true"`
  Description  string        `db:"description"    json:"description,omitempty"     yaml:"description"`
  
  
}

type Project struct {
  CODE         string        `db:"code"           json:"code,omitempty"            yaml:"code"                unique:"true"`
  Name         string        `db:"name"           json:"name,omitempty"            yaml:"name"                unique:"true"`
  ParentProject   string     `db:"parent_task"    json:"parent_task,omitempty"     yaml:"parent_task"`
  Description  string        `db:"description"    json:"description,omitempty"     yaml:"description"`

  StartAt      time.Time     `db:"start_at;default: now()"   json:"start_at"     sql:"default: now()"    gorm:"type:timestamp with time zone"`
  FinishAt     time.Time     `db:"finish_at;default: now()"  json:"finish_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`

  SystemCODE   string        `db:"system_code"    json:"system_code,omitempty"     yaml:"system_code"`

  Milestones   map[string]Milestone           `db:"milestones"    json:"milestones,omitempty"     yaml:"milestones"`
}

type Projects struct {
  a map[string]Project
}

func NewProjects() *Projects {
  return &Projects{
                 a: make(map[string]Project),
               }
}

func (s *Projects) FileExtension() string {
  return ".stage"
}

func (s *Projects) Count() int64 {
  return int64(len(s.a))
}

func (s *Projects) Append(info *Project) {
  info.CODE = strings.ToLower(info.CODE)
  info.ParentProject = strings.ToLower(info.ParentProject)
  s.a[info.CODE] = *info
}

func (s *Projects) GetByCODE(code string) (*Project) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Projects) GetParentProject(code string) string {
  item, ok := s.a[code]
  if ok {
    return item.ParentProject
  }
  return ""
}

func (s *Projects) GetChildProject(code string) string {
  for _, s := range s.a {
    if s.ParentProject == code {
      return s.CODE
    }
  }
  return ""
}

func (s *Projects) GetList() []Project {
  res := make([]Project, 0)
  for _, item := range s.a {
    res = append(res, item)
  }
  return res
}

func (s *Projects) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == s.FileExtension()  {
      if glog.V(2) {
        glog.Infof("LOG: Project file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Project(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Projects: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *Projects) fileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Project

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: ProjectFile(%s): JSON: %v", filename, err)
    return 0
  }
  s.Append(&oTmp)
  return 1
}
