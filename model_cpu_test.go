package arc

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
  
  c := NewCPUs()
  
  c.LoadFromFiles("etc.test")
  
  assert.Nil(t, c.GetByCODE("cpu1"))
  assert.Equal(t, &CPU{Name:"Intel Core i5-4210U @ 1.70GHz", AverageCPUMark:2270, Threads:4, Cores:2}, c.GetByCODE("Intel Core i5-4210U @ 1.70GHz"))
  
  assert.Equal(t, float32(0), c.GetFactor("cpu1", "cpu2"))
  assert.Equal(t, float32(0), c.GetFactor("Intel Core i5-4210U @ 1.70GHz", "cpu2"))
  assert.Equal(t, float32(1), c.GetFactor("Intel Core i5-4210U @ 1.70GHz", "Intel Core i5-4210U @ 1.70GHz"))
  assert.Equal(t, float32(10.054185), c.GetFactor("Intel Core i5-4210U @ 1.70GHz", "AMD Ryzen 7 3700X"))
}
