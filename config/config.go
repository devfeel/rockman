package config

import (
	"errors"
	"github.com/devfeel/rockman/util/file"
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
		RetryLimit            int
		CheckNetInterval      int //the interval time for check net, unit for second
		DataBaseConnectString string
	}

	LoggerSection struct {
		LogPath string
	}

	ClusterSection struct {
		ClusterId             string //cluster Id
		RegistryServer        string //Registry Server
		WatchLeaderRetryLimit int
		QueryResourceInterval int //the interval time for QueryResource, unit for second
	}

	NodeSection struct {
		NodeId   string
		NodeName string
		IsMaster bool
		IsWorker bool
	}

	RpcSection struct {
		OuterHost      string
		OuterPort      string
		RpcHost        string
		RpcPort        string
		RpcProtocol    string //now is json-rpc
		EnableTls      bool
		ServerCertFile string
		ServerKeyFile  string
		ClientCertFile string
		ClientKeyFile  string
	}

	WebUISection struct {
		HttpHost string
		HttpPort string
	}

	RuntimeSection struct {
		IsRun             bool
		LogPath           string
		EnableShellScript bool
	}
)

var CurrentProfile *Profile

func GetConfigPath(file string) string {
	return ConfigPath + file
}

// DefaultProfile return default profile used to full node role
func DefaultProfile() *Profile {
	p := new(Profile)
	p.Global = &GlobalSection{RetryLimit: 5, DataBaseConnectString: "rock:rock@tcp(118.31.32.168:3306)/rockman?charset=utf8&allowOldPasswords=1&loc=Asia%2FShanghai&parseTime=true"}
	p.Node = &NodeSection{NodeId: "a1e97685392845f7b5bbd18f38a10461", IsMaster: true, IsWorker: true}
	p.Rpc = &RpcSection{RpcHost: "", RpcPort: "2398", EnableTls: false, ServerCertFile: "tls/server.crt", ServerKeyFile: "tls/server.key", ClientCertFile: "tls/client.crt", ClientKeyFile: "tls/client.key"}
	p.WebUI = &WebUISection{HttpHost: "", HttpPort: "8080"}
	p.Logger = &LoggerSection{LogPath: "./logs"}
	p.Runtime = &RuntimeSection{IsRun: true, LogPath: "./logs/runtime", EnableShellScript: false}
	p.Cluster = &ClusterSection{RegistryServer: "116.62.16.66:8500", ClusterId: "dev-rock", WatchLeaderRetryLimit: 10, QueryResourceInterval: 60}

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
	if !_file.Exists(realFile) {
		realFile = _file.GetCurrentDirectory() + "/" + configFile
		if !_file.Exists(realFile) {
			realFile = _file.GetCurrentDirectory() + "/config/" + configFile
			if !_file.Exists(realFile) {
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
