package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleConstants(t *testing.T) {
	assert.Equal(t, "user", RoleUser)
	assert.Equal(t, "admin", RoleAdmin)
}
