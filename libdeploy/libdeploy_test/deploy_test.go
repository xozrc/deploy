package libdeploy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
import (
	"github.com/xozrc/deploy/libdeploy"
)

const (
	deployFile = "./machine.yaml"
)

func TestDeploy(t *testing.T) {
	_, err := libdeploy.Deploy(deployFile)
	assert.NoError(t, err, "")
}
