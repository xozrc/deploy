package config

import (
	"io/ioutil"
	"os"
)
import (
	"github.com/ghodss/yaml"
)

type SwarmConfig struct {
	Swarm       bool     `yaml:"Swarm"`
	Master      bool     `yaml:"master"`
	Host        string   `yaml:"host"`
	Image       string   `yaml:"Image"`
	Strategy    string   `yaml:"strategy"`
	Discovery   string   `yaml:"discovery"`
	Expermental bool     `yaml:"Expermental"`
	Opts        []string `yaml:"opts"`
}

type EngineConfig struct {
	StorageDriver    string   `yaml:"aufs"`
	InstallURL       string   `yaml:"installURL"`
	RegistryMirror   []string `yaml:"registryMirror"`
	InsecureRegistry []string `yaml:"insecureRegistry"`
	Labels           []string `yaml:"labels"`
	Env              []string `yaml:"env"`
	Opts             []string `yaml:"opts"`
}

type NodeConfig struct {
	Name             string        `yaml:"name"`
	Driver           string        `yaml:"driver"`
	StorePath        string        `yaml:"storePath"`
	CaCertPath       string        `yaml:"caCertPath"`
	CaPrivateKeyPath string        `yaml:"caPrivateKeyPath"`
	ClientKeyPath    string        `yaml:"clientKeyPath"`
	ClientCertPath   string        `yaml:"clientCertPath"`
	Engine           *EngineConfig `yaml:"engine,omitempty"`
	Swarm            *SwarmConfig  `yaml:"swarm,omitempty"`
}

func NodeConfigFromFile(cfgFile string) (mc *NodeConfig, err error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		return
	}
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	mc = &NodeConfig{}
	err = yaml.Unmarshal(c, mc)
	return
}
