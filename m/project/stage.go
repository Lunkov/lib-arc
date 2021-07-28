package project

import (
  "strings"
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
  
  "github.com/graphql-go/graphql"
  "github.com/SonicRoshan/straf"

  "github.com/Lunkov/lib-arc/gql"
)

type Stage struct {
  CODE         string                   `db:"code"           json:"code,omitempty"            yaml:"code"                unique:"true"`
  Name         string                   `db:"name"           json:"name,omitempty"            yaml:"name"                unique:"true"`
  NextStage    string                   `db:"next_stage"     json:"next_stage,omitempty"      yaml:"next_stage"`
  Description  string                   `db:"description"    json:"description,omitempty"     yaml:"description"`
}

type Stages struct {
  a map[string]Stage
}

func NewStages() *Stages {
  return &Stages{
                 a: make(map[string]Stage),
               }
}

func (s *Stages) FileExtension() string {
  return ".stage"
}

func (s *Stages) Count() int64 {
  return int64(len(s.a))
}

func (s *Stages) Append(info *Stage) {
  info.CODE = strings.ToLower(info.CODE)
  info.NextStage = strings.ToLower(info.NextStage)
  s.a[info.CODE] = *info
}

func (s *Stages) GetByCODE(code string) (*Stage) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Stages) GetNextStage(code string) string {
  item, ok := s.a[code]
  if ok {
    return item.NextStage
  }
  return ""
}

func (s *Stages) GetPrevStage(code string) string {
  for _, s := range s.a {
    if s.NextStage == code {
      return s.CODE
    }
  }
  return ""
}

func (s *Stages) GetList() []Stage {
  res := make([]Stage, 0)
  for _, item := range s.a {
    res = append(res, item)
  }
  return res
}

func (s *Stages) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == s.FileExtension()  {
      if glog.V(2) {
        glog.Infof("LOG: Stage file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Stage(%s)  #%v ", filename, err)
      } else {
        count += s.FileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Stages: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *Stages) FileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Stage

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: StageFile(%s): JSON: %v", filename, err)
    return 0
  }
  s.Append(&oTmp)
  return 1
}

func (s *Stages) InitGQL(g *gql.GQL) {
  StageType, err := straf.GetGraphQLObject(Stage{})
  if err != nil {
    glog.Errorf("ERR: StagesGQL: %s", err)
  }
  
  g.AppendFields("stage", &graphql.Field{
			Type: StageType,
      Args: graphql.FieldConfigArgument{
                "code": &graphql.ArgumentConfig{
                  Description: "code of the Stage",
                  Type:graphql.NewNonNull(graphql.String),
                },
              },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        id := p.Args["code"].(string)
				return s.GetByCODE(id), nil
			},
		})
    
	g.AppendFields("stages", &graphql.Field{
			Type: graphql.NewList(StageType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return s.GetList(), nil
			},
    })
}
