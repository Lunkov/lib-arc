package main

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "gopkg.in/yaml.v2"
)

// Блок диаграммы
// Alt -  Несколько альтернативных фрагментов (alternative); выполняется только тот фрагмент, условие которого истинно
// Opt - Необязательный (optional) фрагмент; выполняется, только если условие истинно. Эквивалентно alt с одной веткой
// Par - Параллельный (parallel); все фрагменты выполняются параллельно
// loop - Цикл (loop); фрагмент может выполняться несколько раз, а защита обозначает тело итерации
// region - Критическая область (critical region); фрагмент может иметь только один поток, выполняющийся за один прием
// Neg - Отрицательный (negative) фрагмент; обозначает неверное взаимодействие

type SequenceBlock struct {
  Type          string            `db:"type"            json:"type"             yaml:"type"`
  Title         string            `db:"title"           json:"title"            yaml:"title"`
  // Links         []SequenceLink    `db:"links"           json:"links"            yaml:"links"`
}

type SequenceLink struct {
  Step     int              `db:"step"      json:"step"      yaml:"step"`
  From     string           `db:"from"      json:"from"      yaml:"from"`
  To       string           `db:"to"        json:"to"        yaml:"to"`
  Type     string           `db:"type"      json:"type"      yaml:"type"`       // request, answer, activate, selfcall
  Comment  string           `db:"comment"   json:"comment"   yaml:"comment"`
  Call     string           `db:"call"      json:"call"      yaml:"call"`
}

type Sequence struct {
  CODE          string            `db:"code"            json:"code"             yaml:"code"`
  UseCase       string            `db:"case"            json:"case"             yaml:"case"`
  Title         string            `db:"title"           json:"title"            yaml:"title"`
  Version       string            `db:"version"         json:"version"          yaml:"version"`
  Description   string            `db:"description"     json:"description"      yaml:"description"`
  Links         []SequenceLink    `db:"links"           json:"links"            yaml:"links"`
  //Blocks        []SequenceBlock   `db:"blocks"          json:"blocks"           yaml:"blocks"`
}

var memSequence = make(map[string]Sequence)

func SequenceCount() int64 {
  return int64(len(memSequence))
}

func SequenceAppend(info *Sequence) {
  memSequence[info.CODE] = *info
}

func GetSequenceByCode(code string) (*Sequence) {
  item, ok := memSequence[code]
  if ok {
    return &item
  }
  return nil
}

func LoadSequencesFromFiles(scanPath string) int {
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false && filepath.Ext(filename) == ".seq"  {
      if glog.V(2) {
        glog.Infof("LOG: Sequence file: %s\n", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.Sequence(%s)  #%v ", filename, err)
      } else {
        count += fileSequenceParse(filename, jsonFile)
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Items: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", scanPath, errScan)
  }

  return count
}

func fileSequenceParse(filename string, jsonFile []byte) int {
  var err error
  var iTmp Sequence

  err = yaml.Unmarshal(jsonFile, &iTmp)
  if err != nil {
    glog.Errorf("ERR: SequenceFile(%s): JSON: %v", filename, err)
  }
  SequenceAppend(&iTmp)

  return 1
}
