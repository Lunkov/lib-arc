package software

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestInteface(t *testing.T) {

  var b Protocol
  
  b.Set(HTTPS)
  
  assert.Equal(t, false, b.Has(HTTP))
  assert.Equal(t, true, b.Has(HTTPS))

  b.Set(HTTP)
  
  assert.Equal(t, true, b.Has(HTTP))
  assert.Equal(t, true, b.Has(HTTPS))
  assert.Equal(t, false, b.Has(TCP))
  assert.Equal(t, false, b.Has(UDP))
  
  b.Clear(HTTPS)

  assert.Equal(t, true, b.Has(HTTP))
  assert.Equal(t, false, b.Has(HTTPS))
  assert.Equal(t, false, b.Has(TCP))
  assert.Equal(t, false, b.Has(UDP))

  b.Zero()
  b.Load([]string{"binary", "proprietary", "TCP"})

  assert.Equal(t, false, b.Has(HTTP))
  assert.Equal(t, false, b.Has(HTTPS))
  assert.Equal(t, true, b.Has(PROPRIETARY))
  assert.Equal(t, true, b.Has(TCP))
  assert.Equal(t, false, b.Has(UDP))
  assert.Equal(t, true, b.Has(BINARY))
  
  assert.Equal(t, []string{"binary", "proprietary", "tcp"}, b.Save())
  
  b.Toggle(HTTPS)
  assert.Equal(t, false, b.Has(HTTP))
  assert.Equal(t, true, b.Has(HTTPS))
  assert.Equal(t, true, b.Has(PROPRIETARY))
  assert.Equal(t, true, b.Has(TCP))
  assert.Equal(t, false, b.Has(UDP))
  assert.Equal(t, true, b.Has(BINARY))
}
