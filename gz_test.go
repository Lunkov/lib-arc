package arc

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestGZ(t *testing.T) {

  s0 := "1234567890"
  s1 := gzdeflate(s0)
  
  assert.Equal(t, "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff2426153\xb7\xb04\x00\x04\x00\x00\xff\xff\xe5\xae\x1d&\n\x00\x00\x00", s1)
  assert.Equal(t, s0, gzinflate(s1))

}

