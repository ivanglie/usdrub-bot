package br

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
	err := SetLogger(log)
	assert.Nil(t, err)
}

func TestSetLogger_Error(t *testing.T) {
	err := SetLogger(nil)
	assert.NotNil(t, err)
}
