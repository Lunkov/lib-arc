package project

import (
  "strings"
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

type Task struct {
  CODE         string                   `db:"code"           json:"code,omitempty"            yaml:"code"                unique:"true"`
  Name         string                   `db:"name"           json:"name,omitempty"            yaml:"name"                unique:"true"`
  ParentTask   string                   `db:"parent_task"    json:"parent_task,omitempty"     yaml:"parent_task"`
  NextTask     string                   `db:"next_task"      json:"next_task,omitempty"       yaml:"next_task"`
  Description  string                   `db:"description"    json:"description,omitempty"     yaml:"description"`
  SystemCODE   string                   `db:"system_code"    json:"system_code,omitempty"     yaml:"system_code"`
}

type Tasks struct {
  a map[string]Task
}

func NewTasks() *Tasks {
  return &Tasks{
                 a: make(map[string]Task),
               }
}

func (s *Tasks) FileExtension() string {
  return ".task"
}

func (s *Tasks) Count() int64 {
  return int64(len(s.a))
}

func (s *Tasks) Append(info *Task) {
  info.CODE = strings.ToLower(info.CODE)
  info.NextTask = strings.ToLower(info.NextTask)
  s.a[info.CODE] = *info
}

func (s *Tasks) GetByCODE(code string) (*Task) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Tasks) GetNextTask(code string) string {
  item, ok := s.a[code]
  if ok {
    return item.NextTask
  }
  return ""
}

func (s *Tasks) GetPrevTask(code string) string {
  for _, s := range s.a {
    if s.NextTask == code {
      return s.CODE
    }
  }
  return ""
}

func (s *Tasks) GetList() []Task {
  res := make([]Task, 0)
  for _, item := range s.a {
    res = append(res, item)
  }
  return res
}

func (s *Tasks) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == s.FileExtension()  {
      if glog.V(2) {
        glog.Infof("LOG: Task file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Task(%s)  #%v ", filename, err)
      } else {
        count += s.fileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Tasks: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *Tasks) fileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Task

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: TaskFile(%s): JSON: %v", filename, err)
    return 0
  }
  s.Append(&oTmp)
  return 1
}
