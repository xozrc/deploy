package libdeploy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)
import (
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/ghodss/yaml"
	"github.com/xozrc/deploy/config"
)

func Deploy(filePath string) (h *host.Host, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	c, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	mc := &config.MachineConfig{}
	if err = yaml.Unmarshal(c, mc); err != nil {
		return
	}

	hostOpts := hostOptsFromMachineConfig(mc)

	driverOpts := driverOptsFromMachineConfig(mc)

	return deploy(mc.Name, mc.Driver, hostOpts, driverOpts)
}

func deploy(hostName, driverName string, hostOpts *host.Options, driverOpts drivers.DriverOptions) (h *host.Host, err error) {

	storePath := mcndirs.GetBaseDir()

	api := libmachine.NewClient(storePath, mcndirs.GetMachineCertDir())
	defer api.Close()

	exists, err := api.Exists(hostName)
	if err != nil {
		return
	}

	//exist
	if exists {
		return nil, mcnerror.ErrHostAlreadyExists{
			Name: h.Name,
		}
	}

	//base driver
	rawDriver, err := json.Marshal(&drivers.BaseDriver{
		MachineName: hostName,
		StorePath:   storePath,
	})

	if err != nil {
		return
	}

	h, err = api.NewHost(driverName, rawDriver)
	if err != nil {
		return
	}

	h.HostOptions.AuthOptions.StorePath = filepath.Join(mcndirs.GetMachineDir(), hostName)
	h.HostOptions.AuthOptions.ServerKeyPath = filepath.Join(mcndirs.GetMachineDir(), hostName, "server.pem")
	h.HostOptions.AuthOptions.ServerCertPath = filepath.Join(mcndirs.GetMachineDir(), hostName, "server-key.pem")

	h.HostOptions = mergeHostOpts(h.HostOptions, hostOpts)

	//set create config
	// if err = h.Driver.SetConfigFromFlags(driverOpts); err != nil {
	// 	return
	// }

	err = api.Create(h)
	if err != nil {
		return
	}

	err = api.Save(h)
	if err != nil {
		return
	}
	return
}

func hostOptsFromMachineConfig(mc *config.MachineConfig) (hostOpts *host.Options) {
	hostOpts = &host.Options{}
	hostOpts.AuthOptions = mc.Auth
	hostOpts.EngineOptions = mc.Engine
	hostOpts.SwarmOptions = mc.Swarm
	return
}

func mergeHostOpts(hostOpts *host.Options, extraHostOpts *host.Options) *host.Options {

	hostOpts.AuthOptions = mergeAuthOpts(hostOpts.AuthOptions, extraHostOpts.AuthOptions)
	hostOpts.EngineOptions = mergeEngineOpts(hostOpts.EngineOptions, extraHostOpts.EngineOptions)

	return hostOpts
}

func mergeAuthOpts(authOpts *auth.Options, extraAuthOpts *auth.Options) *auth.Options {
	if extraAuthOpts == nil {
		return authOpts
	}

	if extraAuthOpts.CaCertPath != "" {
		authOpts.CaCertPath = extraAuthOpts.CaCertPath
	}
	if extraAuthOpts.CaPrivateKeyPath != "" {
		authOpts.CaPrivateKeyPath = extraAuthOpts.CaPrivateKeyPath
	}
	if extraAuthOpts.ClientCertPath != "" {
		authOpts.ClientCertPath = extraAuthOpts.ClientCertPath
	}
	if extraAuthOpts.ClientKeyPath != "" {
		authOpts.ClientKeyPath = extraAuthOpts.ClientKeyPath
	}
	if len(extraAuthOpts.ServerCertSANs) != 0 {
		authOpts.ServerCertSANs = extraAuthOpts.ServerCertSANs
	}

	return authOpts
}

func mergeEngineOpts(engineOpts *engine.Options, extraEngineOpts *engine.Options) *engine.Options {
	if extraEngineOpts == nil {
		return engineOpts
	}
	if len(extraEngineOpts.ArbitraryFlags) != 0 {
		engineOpts.ArbitraryFlags = extraEngineOpts.ArbitraryFlags
	}

	if len(extraEngineOpts.DNS) != 0 {
		engineOpts.DNS = extraEngineOpts.DNS
	}

	if len(extraEngineOpts.Env) != 0 {
		engineOpts.Env = extraEngineOpts.Env
	}

	if extraEngineOpts.GraphDir != "" {
		engineOpts.GraphDir = extraEngineOpts.GraphDir
	}

	if len(extraEngineOpts.InsecureRegistry) != 0 {
		engineOpts.InsecureRegistry = extraEngineOpts.InsecureRegistry
	}

	if extraEngineOpts.InstallURL != "" {
		engineOpts.InstallURL = extraEngineOpts.InstallURL
	}

	engineOpts.Ipv6 = extraEngineOpts.Ipv6

	if len(extraEngineOpts.Labels) != 0 {
		engineOpts.Labels = extraEngineOpts.Labels
	}

	if extraEngineOpts.LogLevel != "" {
		engineOpts.LogLevel = extraEngineOpts.LogLevel
	}
	if len(extraEngineOpts.RegistryMirror) != 0 {
		engineOpts.RegistryMirror = extraEngineOpts.RegistryMirror
	}

	engineOpts.SelinuxEnabled = extraEngineOpts.SelinuxEnabled

	if extraEngineOpts.StorageDriver != "" {
		engineOpts.StorageDriver = extraEngineOpts.StorageDriver
	}
	engineOpts.TLSVerify = extraEngineOpts.TLSVerify
	return engineOpts

}

func driverOptsFromMachineConfig(mc *config.MachineConfig) drivers.DriverOptions {
	return nil
}
