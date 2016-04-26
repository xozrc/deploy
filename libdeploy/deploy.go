package libdeploy

import (
	"encoding/json"
	"path/filepath"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/swarm"
	"github.com/xozrc/deploy/config"
)

func Deploy(mc *config.MachineConfig) (h *host.Host, err error) {
	storePath := ""
	if mc.Auth != nil {
		storePath = mc.Auth.StorePath
	}

	hostOpts := hostOptsFromMachineConfig(mc)
	driverOpts := driverOptsFromMachineConfig(mc)

	return deploy(storePath, mc.Name, mc.Driver, hostOpts, driverOpts)
}

func Undeploy(mc *config.MachineConfig) {
	storePath := mc.Auth.StorePath
	api := libmachine.NewClient(storePath, getMachineCertDir(storePath))

	hostName := mc.Name

	defer api.Close()
	h, err := api.Load(hostName)
	if err != nil {
		goto RemoveLocal
	}
	h.Driver.Remove()
RemoveLocal:
	api.Remove(hostName)
}

func deploy(storePath string, hostName, driverName string, hostOpts *host.Options, driverOpts drivers.DriverOptions) (h *host.Host, err error) {

	api := libmachine.NewClient(storePath, getMachineCertDir(storePath))

	defer api.Close()

	exists, err := api.Exists(hostName)
	if err != nil {
		return
	}

	//exist
	if exists {
		h, err = api.Load(hostName)

		return
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

	h.HostOptions.AuthOptions.StorePath = filepath.Join(getMachineDir(storePath), hostName)
	h.HostOptions.AuthOptions.ServerKeyPath = filepath.Join(getMachineDir(storePath), hostName, "server.pem")
	h.HostOptions.AuthOptions.ServerCertPath = filepath.Join(getMachineDir(storePath), hostName, "server-key.pem")

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
	hostOpts.SwarmOptions = mergeSwarmOpts(hostOpts.SwarmOptions, extraHostOpts.SwarmOptions)
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

	return authOpts
}

func mergeEngineOpts(engineOpts *engine.Options, extraEngineOpts *engine.Options) *engine.Options {
	if extraEngineOpts == nil {
		return engineOpts
	}
	if len(extraEngineOpts.ArbitraryFlags) != 0 {
		engineOpts.ArbitraryFlags = extraEngineOpts.ArbitraryFlags
	}

	if len(extraEngineOpts.Env) != 0 {
		engineOpts.Env = extraEngineOpts.Env
	}

	if len(extraEngineOpts.InsecureRegistry) != 0 {
		engineOpts.InsecureRegistry = extraEngineOpts.InsecureRegistry
	}

	if extraEngineOpts.InstallURL != "" {
		engineOpts.InstallURL = extraEngineOpts.InstallURL
	}

	if len(extraEngineOpts.Labels) != 0 {
		engineOpts.Labels = extraEngineOpts.Labels
	}

	if len(extraEngineOpts.RegistryMirror) != 0 {
		engineOpts.RegistryMirror = extraEngineOpts.RegistryMirror
	}

	return engineOpts

}

func mergeSwarmOpts(swarmOpts *swarm.Options, extraSwarmOpts *swarm.Options) *swarm.Options {
	if extraSwarmOpts == nil {
		return swarmOpts
	}

	if extraSwarmOpts.Host == "" {
		extraSwarmOpts.Host = swarmOpts.Host
	}

	if extraSwarmOpts.Image == "" {
		extraSwarmOpts.Image = swarmOpts.Image
	}

	if extraSwarmOpts.Strategy == "" {
		extraSwarmOpts.Strategy = swarmOpts.Strategy
	}
	return extraSwarmOpts

}

func driverOptsFromMachineConfig(mc *config.MachineConfig) drivers.DriverOptions {
	return nil
}

func getMachineCertDir(baseDir string) string {
	return filepath.Join(baseDir, "certs")
}
func getMachineDir(baseDir string) string {
	return filepath.Join(baseDir, "machines")
}
