package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingHostOrUser(t *testing.T) {
	plugin := Plugin{}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, missingHostOrUser, err.Error())
}

func TestMissingKeyOrPassword(t *testing.T) {
	plugin := Plugin{
		Config{
			Host: []string{"localhost"},
			User: "ubuntu",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, missingPasswordOrKey, err.Error())
}
