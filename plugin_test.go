package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingKeyOrPassword(t *testing.T) {
	plugin := Plugin{
		Config{
			Host: []string{"localhost"},
			User: "ubuntu",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}
