package arc

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {

  s := NewRoles()
  s.LoadFromFiles("etc.test")

  assert.Equal(t, int64(2), s.Count())
  assert.Equal(t, (*Role)(nil), s.GetByCODE("admin1"))
  assert.Equal(t, &Role{CODE:"admin", Type:"user", Name:"Администратор", Description:""}, s.GetByCODE("admin"))

}

