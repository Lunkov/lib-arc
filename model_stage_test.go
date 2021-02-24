package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestStage(t *testing.T) {
  
  s := NewStages()
  s.LoadFromFiles("etc.test")
  
  assert.Equal(t, int64(7), s.Count())
  assert.Equal(t, (*Stage)(nil), s.GetByCODE("stage1"))
  assert.Equal(t, &Stage{CODE:"design", Name:"Design", NextStage:"develop", Description:""}, s.GetByCODE("design"))
  
  assert.Equal(t, "pilot", s.GetNextStage("test"))
  assert.Equal(t, "develop", s.GetPrevStage("test"))
}

