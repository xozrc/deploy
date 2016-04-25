package config

import (
	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/swarm"
)

// type AuthConfig struct {
// 	CertDir              string   `yaml:"certDir,omitempty"`
// 	CaCertPath           string   `yaml:"caCertPath,omitempty"`
// 	CaPrivateKeyPath     string   `yaml:"caPrivateKeyPath,omitempty"`
// 	CaCertRemotePath     string   `yaml:"caCertRemotePath,omitempty"`
// 	ServerCertPath       string   `yaml:"serverCertPath,omitempty"`
// 	ServerKeyPath        string   `yaml:"serverKeyPath,omitempty"`
// 	ClientKeyPath        string   `yaml:"clientKeyPath,omitempty"`
// 	ServerCertRemotePath string   `yaml:"serverCertRemotePath,omitempty"`
// 	ServerKeyRemotePath  string   `yaml:"serverKeyRemotePath,omitempty"`
// 	ClientCertPath       string   `yaml:"clientCertPath,omitempty"`
// 	ServerCertSANs       []string `yaml:"serverCertSANs,omitempty"`
// 	StorePath            string   `yaml:"storePath,omitempty"`
// }

// type EngineConfig struct {
// 	ArbitraryFlags   []string `yaml:"engineOpt,omitempty"`
// 	DNS              []string `yaml:"dns,omitempty"`
// 	GraphDir         string   `yaml:"graphDir,omitempty"`
// 	Env              []string `yaml:"env,omitempty"`
// 	Ipv6             bool     `yaml:"ipv6,omitempty"`
// 	InsecureRegistry []string `yaml:"insecureRegistry,omitempty"`
// 	Labels           []string `yaml:"labels,omitempty"`
// 	LogLevel         string   `yaml:"logLevel,omitempty"`
// 	StorageDriver    string   `yaml:"storageDriver,omitempty"`
// 	SelinuxEnabled   bool     `yaml:"selinuxEnabled,omitempty"`
// 	TLSVerify        bool     `yaml:"tLSVerify,omitempty"`
// 	RegistryMirror   []string `yaml:"registryMirror,omitempty"`
// 	InstallURL       string   `yaml:"installURL,omitempty"`
// }

// type SwarmConfig struct {
// 	IsSwarm        bool     `yaml:"isSwarm,omitempty"`
// 	Address        string   `yaml:"address,omitempty"`
// 	Discovery      string   `yaml:"discovery,omitempty"`
// 	Master         bool     `yaml:"master,omitempty"`
// 	Host           string   `yaml:"host,omitempty"`
// 	Image          string   `yaml:"image,omitempty"`
// 	Strategy       string   `yaml:"strategy,omitempty"`
// 	Heartbeat      int      `yaml:"heartbeat,omitempty"`
// 	Overcommit     float64  `yaml:"overcommit,omitempty"`
// 	ArbitraryFlags []string `yaml:"arbitraryFlags,omitempty"`
// 	Env            []string `yaml:"env,omitempty"`
// 	IsExperimental bool     `yaml:"isExperimental,omitempty"`
// }

type MachineConfig struct {
	Name   string          `yaml:"Name"`
	Driver string          `yaml:"Driver"`
	Auth   *auth.Options   `yaml:"Auth"`
	Engine *engine.Options `yaml:"Engine"`
	Swarm  *swarm.Options  `yaml:"Swarm"`
}
