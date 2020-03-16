package config

import (
	"errors"
	"github.com/devfeel/rockman/src/util/file"
	"github.com/devfeel/rockman/src/util/uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const ConfigPath = "./conf/"

type (
	Profile struct {
		Global  *GlobalSection
		Node    *NodeSection
		Rpc     *RpcSection
		WebUI   *WebUISection
		Runtime *RuntimeSection
		Logger  *LoggerSection
		Cluster *ClusterSection
	}

	GlobalSection struct {
		DataBaseConnectString string
	}

	LoggerSection struct {
		LogPath string
	}

	ClusterSection struct {
		ClusterId      string //cluster Id
		RegistryServer string //Registry Server
	}

	NodeSection struct {
		NodeId   string
		NodeName string
		IsMaster bool
		IsWorker bool
	}

	RpcSection struct {
		RpcHost     string
		RpcPort     string
		RpcProtocol string //now is json-rpc
	}

	WebUISection struct {
		HttpHost string
		HttpPort string
	}

	RuntimeSection struct {
		IsRun   bool
		LogPath string
	}
)

var CurrentProfile *Profile

func GetConfigPath(file string) string {
	return ConfigPath + file
}

// SingleNodeProfile return default profile used to single node
func SingleNodeProfile() *Profile {
	p := new(Profile)
	p.Global = &GlobalSection{DataBaseConnectString: "rock:rock@tcp(118.31.32.168:3306)/rockman?charset=utf8&allowOldPasswords=1"}
	p.Node = &NodeSection{NodeId: uuid.NewV4().String32(), IsMaster: true, IsWorker: true}
	p.Rpc = &RpcSection{RpcHost: "127.0.0.1", RpcPort: "2398"}
	p.WebUI = &WebUISection{HttpHost: "127.0.0.1", HttpPort: "8080"}
	p.Logger = &LoggerSection{LogPath: "./logs"}
	p.Runtime = &RuntimeSection{IsRun: true, LogPath: "./logs/runtime"}
	p.Cluster = &ClusterSection{RegistryServer: "116.62.16.66:8500", ClusterId: "rock"}

	CurrentProfile = p
	return p
}

// LoadConfig load config from yaml file
func LoadConfig(configFile string) (*Profile, error) {
	// Validity check
	// 1. Try read as absolute path
	// 2. Try the current working directory
	// 3. Try $PWD/config
	// fixed for issue #15 config file path
	realFile := configFile
	if !file.Exist(realFile) {
		realFile = file.GetCurrentDirectory() + "/" + configFile
		if !file.Exist(realFile) {
			realFile = file.GetCurrentDirectory() + "/config/" + configFile
			if !file.Exist(realFile) {
				return nil, errors.New("no exists config file => " + configFile)
			}
		}
	}

	profile, err := loadYamlConfig(realFile)
	if err == nil {
		err = validateProfile(profile)
	}
	CurrentProfile = profile
	return profile, err
}

// loadYamlConfig load config profile from yaml file
func loadYamlConfig(configFile string) (*Profile, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New("RockMan:Config:initYamlConfig config file [" + configFile + "] cannot be parsed - " + err.Error())
	}

	var profile *Profile
	err = yaml.Unmarshal(content, &profile)
	if err != nil {
		return nil, errors.New("RockMan:Config:initYamlConfig config file [" + configFile + "] cannot be parsed - " + err.Error())
	}
	return profile, nil
}

func validateProfile(profile *Profile) error {
	if profile.Node.NodeId == "" {
		return errors.New("RockMan start error: " + "NodeId is null")
	}
	return nil
}
