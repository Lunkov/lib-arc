package main

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCPU(t *testing.T) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  // flag.Set("v", "9")
  flag.Parse()
  
  LoadCPUsFromFiles("etc.test/cpu")
  
  assert.Nil(t, GetCPUByCODE("cpu1"))
  assert.Equal(t, &CPU{Name:"Intel Core i5-4210U @ 1.70GHz", AverageCPUMark:2270, Threads:4, Cores:2}, GetCPUByCODE("Intel Core i5-4210U @ 1.70GHz"))
  
  assert.Equal(t, float32(0), GetCPUFactor("cpu1", "cpu2"))
  assert.Equal(t, float32(0), GetCPUFactor("Intel Core i5-4210U @ 1.70GHz", "cpu2"))
  assert.Equal(t, float32(1), GetCPUFactor("Intel Core i5-4210U @ 1.70GHz", "Intel Core i5-4210U @ 1.70GHz"))
  assert.Equal(t, float32(10.054185), GetCPUFactor("Intel Core i5-4210U @ 1.70GHz", "AMD Ryzen 7 3700X"))
}
