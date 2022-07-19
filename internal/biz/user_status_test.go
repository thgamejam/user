package biz

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserStatusExist(t *testing.T) {
	s := NewUserStatus(0)
	assert.False(t, s.IsExist())

	s.SetExist(true)
	assert.True(t, s.IsExist())

	s.SetExist(false)
	assert.False(t, s.IsExist())
}

func TestUserStatusBan(t *testing.T) {
	s := NewUserStatus(0)
	assert.False(t, s.IsBan())

	s.SetBan(true)
	assert.True(t, s.IsBan())

	s.SetBan(false)
	assert.False(t, s.IsBan())
}

func TestUserStatusMix(t *testing.T) {
	s := NewUserStatus(0)
	assert.False(t, s.IsExist())
	assert.False(t, s.IsBan())

	s.SetExist(true)
	assert.True(t, s.IsExist())
	assert.False(t, s.IsBan())

	s.SetBan(true)
	assert.True(t, s.IsExist())
	assert.True(t, s.IsBan())

	s.SetExist(false)
	assert.False(t, s.IsExist())
	assert.True(t, s.IsBan())

	s.SetExist(true)
	s.SetBan(false)
	assert.True(t, s.IsExist())
	assert.False(t, s.IsBan())
}
