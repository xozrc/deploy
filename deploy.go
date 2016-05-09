package deploy

import (
	"encoding/json"
	"path/filepath"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/state"
	"github.com/docker/machine/libmachine/swarm"
	"github.com/xozrc/deploy/config"
)

func Undeploy(nc *config.NodeConfig) {
	storePath := nc.StorePath
	api := libmachine.NewClient(storePath, getMachineCertDir(storePath))

	hostName := nc.Name

	defer api.Close()
	h, err := api.Load(hostName)
	if err != nil {
		goto RemoveLocal
	}
	h.Driver.Remove()
RemoveLocal:
	api.Remove(hostName)
}

func Deploy(nc *config.NodeConfig) (h *host.Host, err error) {

	hostName := nc.Name
	driverName := nc.Driver
	storePath := nc.StorePath

	api := libmachine.NewClient(storePath, getMachineCertDir(storePath))
	defer api.Close()

	exists, err := api.Exists(hostName)
	if err != nil {
		return
	}

	//exist
	if exists {
		h, err = api.Load(hostName)

		if err != nil {
			return nil, err
		}

		ts, err := h.Driver.GetState()
		//already running
		if ts == state.Running {
			return h, nil
		}

		err = h.Restart()
		if err != nil {
			return nil, err
		}

		//check swarm

		//fix:test certs
		err = h.ConfigureAuth()

		return h, err
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

	//config auth
	h.HostOptions.AuthOptions.StorePath = filepath.Join(getMachineDir(storePath), hostName)
	h.HostOptions.AuthOptions.ServerKeyPath = filepath.Join(getMachineDir(storePath), hostName, "server.pem")
	h.HostOptions.AuthOptions.ServerCertPath = filepath.Join(getMachineDir(storePath), hostName, "server-key.pem")
	if nc.CaCertPath != "" {
		h.HostOptions.AuthOptions.CaCertPath = nc.CaCertPath
	}
	if nc.CaCertPath != "" {
		h.HostOptions.AuthOptions.CaPrivateKeyPath = nc.CaPrivateKeyPath
	}
	if nc.CaCertPath != "" {
		h.HostOptions.AuthOptions.ClientCertPath = nc.ClientCertPath
	}
	if nc.CaCertPath != "" {
		h.HostOptions.AuthOptions.ClientKeyPath = nc.ClientKeyPath
	}

	//config engine
	if nc.Engine != nil {
		h.HostOptions.EngineOptions = mergeEngineOpts(h.HostOptions.EngineOptions, nc.Engine)
	}

	//config swarm
	if nc.Swarm != nil {
		h.HostOptions.SwarmOptions = mergeSwarmOpts(h.HostOptions.SwarmOptions, nc.Swarm)
	}

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

func mergeEngineOpts(engineOpts *engine.Options, engineCfg *config.EngineConfig) *engine.Options {

	if len(engineCfg.Opts) != 0 {
		engineOpts.ArbitraryFlags = engineCfg.Opts
	}

	if len(engineCfg.Env) != 0 {
		engineOpts.Env = engineCfg.Env
	}

	if len(engineCfg.InsecureRegistry) != 0 {
		engineOpts.InsecureRegistry = engineCfg.InsecureRegistry
	}

	if engineCfg.InstallURL != "" {
		engineOpts.InstallURL = engineCfg.InstallURL
	}

	if len(engineCfg.Labels) != 0 {
		engineOpts.Labels = engineCfg.Labels
	}

	if len(engineCfg.RegistryMirror) != 0 {
		engineOpts.RegistryMirror = engineCfg.RegistryMirror
	}

	if engineCfg.StorageDriver != "" {
		engineOpts.StorageDriver = engineCfg.StorageDriver
	}

	return engineOpts

}

func mergeSwarmOpts(swarmOpts *swarm.Options, swarmCfg *config.SwarmConfig) *swarm.Options {

	swarmOpts.IsSwarm = swarmCfg.Swarm
	swarmOpts.Master = swarmCfg.Master
	swarmOpts.IsExperimental = swarmCfg.Expermental
	if swarmCfg.Host != "" {
		swarmOpts.Host = swarmCfg.Host
	}

	if swarmCfg.Image != "" {
		swarmOpts.Image = swarmCfg.Image
	}

	if swarmCfg.Strategy != "" {
		swarmOpts.Strategy = swarmCfg.Strategy
	}

	if swarmCfg.Discovery != "" {
		swarmOpts.Discovery = swarmCfg.Discovery
	}

	if len(swarmCfg.Opts) == 0 {
		swarmOpts.ArbitraryFlags = swarmCfg.Opts
	}

	return swarmOpts

}

func driverOptsFromMachineConfig(mc *config.NodeConfig) drivers.DriverOptions {
	return nil
}

func getMachineCertDir(baseDir string) string {
	return filepath.Join(baseDir, "certs")
}
func getMachineDir(baseDir string) string {
	return filepath.Join(baseDir, "machines")
}
