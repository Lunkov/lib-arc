package software

import (
  "os"
  "io/ioutil"
  "strings"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"

  "github.com/graphql-go/graphql"
  "github.com/SonicRoshan/straf"

  "github.com/Lunkov/lib-arc/gql"
)

type System struct {
  CODE         string    `db:"code"           json:"code,omitempty"            yaml:"code"`
  ParentCODE   string    `db:"parent_code"    json:"parent_code,omitempty"     yaml:"parent_code"`
  Type         string    `db:"type"           json:"type,omitempty"            yaml:"type"`
  Name         string    `db:"name"           json:"name,omitempty"            yaml:"name"`
  Disabled     bool      `db:"disabled"       json:"disabled,omitempty"        yaml:"disabled"`
}


type Systems struct {
  a map[string]System
}

func NewSystems() *Systems {
  return &Systems{
                 a: make(map[string]System),
               }
}

func (s *Systems) FileExtension() string {
  return ".system"
}

func (s *Systems) Count() int64 {
  return int64(len(s.a))
}

func (s *Systems) Append(info *System) {
  info.CODE = strings.ToLower(info.CODE)
  s.a[info.CODE] = *info
}

func (s *Systems) GetByCODE(code string) (*System) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Systems) GetList() []System {
  res := make([]System, 0)
  for _, item := range s.a {
    res = append(res, item)
  }
  return res
}

func (s *Systems) LoadFromFiles(scanPath string) int {
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

func (s *Systems) FileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp System

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

func (s *Systems) InitGQL(g *gql.GQL) {
  SystemType, err := straf.GetGraphQLObject(System{})
  if err != nil {
    glog.Errorf("ERR: SystemGQL: %s", err)
  }

  g.AppendFields("system", &graphql.Field{
			Type: SystemType,
      Args: graphql.FieldConfigArgument{
                "code": &graphql.ArgumentConfig{
                  Description: "code of the Service",
                  Type:graphql.NewNonNull(graphql.String),
                },
              },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        id := p.Args["code"].(string)
				return s.GetByCODE(id), nil
			},
		})
    
	g.AppendFields("systems", &graphql.Field{
			Type: graphql.NewList(SystemType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return s.GetList(), nil
			},
    })
}
