package main

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  // flag.Set("v", "9")
  flag.Parse()
  
  LoadServicesFromFiles("etc.test")
  
  assert.Equal(t, int64(5), ServiceCount())
}
