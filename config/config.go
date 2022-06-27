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
	ProtocolsData []ProtocolsData `yaml:"protocols_data"`
}

type ProtocolsData struct {
	ProtocolName    string        `yaml:"protocol_name"`
	ChainName       string        `yaml:"chain_name"`
	ContractAddress string        `yaml:"contractAddress"`
	WalletAddress   string        `yaml:"walletAddress"`
	RPC             string        `yaml:"rpc"`
	Stake           StakeConfig   `yaml:"stake"`
	Unstake         UnstakeConfig `yaml:"unstake"`
}

type StakeConfig struct {
	Action           string `yaml:"action"`
	AmountRequired   string `yaml:"amount_required"`
	DurationRequired string `yaml:"duration_required"`
}

type UnstakeConfig struct {
	Action           string `yaml:"action"`
	AmountRequired   string `yaml:"amount_required"`
	DurationRequired string `yaml:"duration_required"`
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
	// env := os.Getenv("ENVIRONMENT")
	// switch env {
	// case "DEVELOPMENT":
	// 	protocolName = "default_development"
	// default:
	// 	err := errors.New("invalid environment")
	// 	log.Fatalf("Error in Setting Environment : %v", err)
	// }

	// //For deployment config
	// viper.SetConfigName(protocolName)
	// viper.AddConfigPath("./config")
	// viper.SetConfigType("yml")
	// err := viper.MergeInConfig()
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// err = viper.Unmarshal(&protocolConfiguration)
	// if err != nil {
	// 	log.Fatalf("Unable to decode into struct, %v", err)
	// }
	return protocolConfiguration
}
