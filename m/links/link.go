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

type Link struct {
  CODE         string                   `db:"code"           json:"code,omitempty"            yaml:"code"             unique:"true"`
  Layer        string                   `db:"layer"          json:"layer,omitempty"           yaml:"layer"`
  Name         string                   `db:"name"           json:"name,omitempty"            yaml:"name"`
  Description  string                   `db:"description"    json:"description,omitempty"     yaml:"description"`
  NodeCode     string                   `db:"node_code"      json:"node_code,omitempty"       yaml:"node_code"`
  SystemCode   string                   `db:"system_code"    json:"system_code,omitempty"     yaml:"system_code"`
  Disabled     bool                     `db:"disabled"       json:"disabled,omitempty"        yaml:"disabled"`
}

type Links struct {
  a map[string]Link
}

func NewLinks() *Links {
  return &Links{
                 a: make(map[string]Link),
               }
}

func (s *Links) FileExtension() string {
  return ".service"
}

func (s *Links) Count() int64 {
  return int64(len(s.a))
}

func (s *Links) Append(info *Link) {
  info.CODE = strings.ToLower(info.CODE)
  s.a[info.CODE] = *info
}

func (s *Links) GetByCODE(code string) (*Link) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Links) GetList() []Link {
  res := make([]Link, 0)
  for _, item := range s.a {
    res = append(res, item)
  }
  return res
}

func (s *Links) LoadFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == s.FileExtension()  {
      if glog.V(2) {
        glog.Infof("LOG: Link file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Link(%s)  #%v ", filename, err)
      } else {
        count += s.FileParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Links: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func (s *Links) FileParse(filename string, jsonFile []byte) int {
  var err error
  var oTmp Link

  err = yaml.Unmarshal(jsonFile, &oTmp)
  if err != nil {
    glog.Errorf("ERR: LinkFile(%s): JSON: %v", filename, err)
  }
  if !oTmp.Disabled {
    s.Append(&oTmp)
    return 1
  }

  return 0
}

func ExportInterfacesFromSwagger() {
}


func (s *Links) InitGQL(g *gql.GQL) {
  LinkType, err := straf.GetGraphQLObject(Link{})
  if err != nil {
    glog.Errorf("ERR: LinkGQL: %s", err)
  }

  g.AppendFields("link", &graphql.Field{
			Type: LinkType,
      Args: graphql.FieldConfigArgument{
                "code": &graphql.ArgumentConfig{
                           Description: "code of the Link",
                           Type: graphql.NewNonNull(graphql.String),
                },
              },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        id := p.Args["code"].(string)
				return s.GetByCODE(id), nil
			},
		})
    
	g.AppendFields("links", &graphql.Field{
			Type: graphql.NewList(LinkType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return s.GetList(), nil
			},
    })
}
