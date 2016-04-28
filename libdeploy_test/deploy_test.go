package libdeploy_test

import "testing"

import (
	"github.com/stretchr/testify/assert"
	"github.com/xozrc/deploy"
	"github.com/xozrc/deploy/config"
)

const (
	deployFile = "./machine.yaml"
)

func TestDeploy(t *testing.T) {
	mc, err := config.NodeConfigFromFile(deployFile)

	if !assert.NoError(t, err, "") {
		return
	}
	_, err = deploy.Deploy(mc)

	assert.NoError(t, err, "")
}
