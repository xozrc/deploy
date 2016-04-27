package libdeploy_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xozrc/deploy/config"
	"github.com/xozrc/deploy/libdeploy"
)

const (
	deployFile = "./machine.yaml"
)

func TestDeploy(t *testing.T) {
	mc, err := config.MachineConfigFromFile(deployFile)
	fmt.Printf("%v,%v\n", mc.Auth, mc.Swarm)
	if !assert.NoError(t, err, "") {
		return
	}
	_, err = libdeploy.Deploy(mc)
	defer libdeploy.Undeploy(mc)
	assert.NoError(t, err, "")

}
