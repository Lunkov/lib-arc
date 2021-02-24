package main

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestUseCase(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  uc := NewUseCaseSet()
  
  uc.LoadFromFiles("etc.test")
  
  assert.Equal(t, int64(2), uc.Count())
  assert.Nil(t, uc.GetByCODE("UC-0001-00"))
  assert.Equal(t, &UseCase{CODE:"UC-0001", Type:"", Path:"etc.test/local", Date:"", Name:"Upload images", Disabled:false, URL:"", Services:[]string(nil), Systems:[]string(nil), Sequences:[]string(nil), Tasks:[]string(nil), ReadMe:""}, uc.GetByCODE("UC-0001"))
  
}
