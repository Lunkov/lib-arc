package software

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

type Deploy struct {
  CODE         string                   `db:"code"           json:"code,omitempty"            yaml:"code"             unique:"true"`
  Layer        string                   `db:"layer"          json:"layer,omitempty"           yaml:"layer"`
  Name         string                   `db:"name"           json:"name,omitempty"            yaml:"name"`
  Description  string                   `db:"description"    json:"description,omitempty"     yaml:"description"`
  NodeCode     string                   `db:"node_code"      json:"node_code,omitempty"       yaml:"node_code"`
  SystemCode   string                   `db:"system_code"    json:"system_code,omitempty"     yaml:"system_code"`
  Disabled     bool                     `db:"disabled"       json:"disabled,omitempty"        yaml:"disabled"`
}

type Deploys struct {
  a map[string]Deploy
}

func NewDeploys() *Deploys {
  return &Deploys{
                 a: make(map[string]Deploy),
               }
}

func (s *Deploys) FileExtension() string {
  return ".service"
}

func (s *Deploys) Count() int64 {
  return int64(len(s.a))
}

func (s *Deploys) Append(info *Deploy) {
  info.CODE = strings.ToLower(info.CODE)
  s.a[info.CODE] = *info
}

func (s *Deploys) GetByCODE(code string) (*Deploy) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Deploys) GetList() []Deploy {
  res := make([]Deploy, 0)
  for _, item := range s.a {
    res = append(res, item)
  }
  return res
}

func (s *Deploys) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == s.FileExtension()  {
      if glog.V(2) {
        glog.Infof("LOG: Deploy file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Deploy(%s)  #%v ", filename, err)
      } else {
        count += s.FileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Deploys: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *Deploys) FileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Deploy

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: DeployFile(%s): JSON: %v", filename, err)
  }
  if !oTmp.Disabled {
    s.Append(&oTmp)
    return 1
  }

  return 0
}

func ExportInterfacesFromSwagger() {
}


func (s *Deploys) InitGQL(g *gql.GQL) {
  DeployType, err := straf.GetGraphQLObject(Deploy{})
  if err != nil {
    glog.Errorf("ERR: DeployGQL: %s", err)
  }

  g.AppendFields("deploy", &graphql.Field{
			Type: DeployType,
      Args: graphql.FieldConfigArgument{
                "code": &graphql.ArgumentConfig{
                           Description: "code of the Deploy",
                           Type: graphql.NewNonNull(graphql.String),
                },
              },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        id := p.Args["code"].(string)
				return s.GetByCODE(id), nil
			},
		})
    
	g.AppendFields("deploys", &graphql.Field{
			Type: graphql.NewList(DeployType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return s.GetList(), nil
			},
    })
}
