package arc

import (
  "strings"
)

// Protocols.
type Protocol uint64

const (
  HTTP Protocol = 1 << iota
  HTTPS
  RESTAPI
  GRPC
  BINARY
  PROPRIETARY
  TCP
  UDP
)

var protocols = map[string]Protocol{
  "http":         HTTP,
  "https":        HTTPS,
  "restapi":      RESTAPI,
  "grps":         GRPC,
  "binary":       BINARY,
  "proprietary":  PROPRIETARY,
  "tcp":          TCP,
  "udp":          UDP,
}

// Interface
type Interface struct {
  CODE           string                   `db:"code"           json:"code,omitempty"            yaml:"code"             unique:"true"`
  Type           string                   `db:"type"           json:"type,omitempty"            yaml:"type"`
  Protocol       Protocol                 `db:"protocol"       json:"protocol,omitempty"        yaml:"protocol"`
  ProtocolTags []string                   `db:"protocol_tags"  json:"protocol_tags,omitempty"   yaml:"protocol_tags"`
  Name           string                   `db:"name"           json:"name,omitempty"            yaml:"name"`
  Description    string                   `db:"description"    json:"description"               yaml:"description"`
  Port           uint                     `db:"port"           json:"port,omitempty"            yaml:"port"`
  
  Input          DataSet                  `db:"input"          json:"input,omitempty"           yaml:"input"          gorm:"column:input;type:jsonb;"`
  Output         DataSet                  `db:"output"         json:"output,omitempty"          yaml:"output"         gorm:"column:output;type:jsonb;"`
}

func (b *Protocol) Zero()                    { (*b) = 0 }
func (b *Protocol) Set(flag Protocol)        { (*b)|= flag }
func (b *Protocol) Clear(flag Protocol)      { (*b) = (*b) &^ flag }
func (b *Protocol) Toggle(flag Protocol)     { (*b)^= flag }
func (b *Protocol) Has(flag Protocol) bool   { return (*b)&flag != 0 }

func (b *Protocol) Load(flags []string) {
  for _, s := range flags {
    b.Set(protocols[strings.ToLower(s)])
  }
}

func (b *Protocol) Save() []string {
  res := make([]string, 0)
  for i, v := range protocols {
    if b.Has(v) {
      res = append(res, i)
    }
  }
  return res
}
