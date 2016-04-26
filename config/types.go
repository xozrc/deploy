package config

import (
	"io/ioutil"
	"os"

	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/swarm"
	"github.com/ghodss/yaml"
)

type MachineConfig struct {
	Name   string
	Driver string
	Auth   *auth.Options
	Engine *engine.Options
	Swarm  *swarm.Options
}

func MachineConfigFromFile(cfgFile string) (mc *MachineConfig, err error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		return
	}
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	mc = &MachineConfig{}
	err = yaml.Unmarshal(c, mc)

	return

}
