package config

import (
	"log"
	"runtime"

	"github.com/spf13/viper"
)

type ProtocolConfig struct {
	Protocols Protocols
}

type Protocols struct {
	ProtocolData []ProtocolData `yaml:"protocols_data"`
}

type ProtocolData struct {
	ProtocolName    string        `yaml:"protocol_name"`
	ChainName       string        `yaml:"chain_name"`
	ContractAddress string        `yaml:"contractAddress"`
	WalletAddress   string        `yaml:"walletAddress"`
	RPC             string        `yaml:"rpc"`
	Stake           StakeConfig   `yaml:"stake"`
	Unstake         UnstakeConfig `yaml:"unstake"`
}

type StakeConfig struct {
	RequiredArray [3]string `yaml:"required_array"`
}

type UnstakeConfig struct {
	RequiredArray [3]string `yaml:"required_array"`
}

func LoadProtocol() ProtocolConfig {
	var protocolConfiguration ProtocolConfig
	// var protocolName string
	//For Local Config
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		viper.SetConfigName("default_dev")
		// Set the path to look for the configurations file
		viper.AddConfigPath("./config")
		// Enable VIPER to read Environment Variables
		viper.AutomaticEnv()
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err.Error())
		}
		err := viper.Unmarshal(&protocolConfiguration)
		if err != nil {
			log.Fatalf("Unable to decode into struct, %v", err)
		}
		return protocolConfiguration
	}
	return protocolConfiguration
}
