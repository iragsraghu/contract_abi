package mapping

import (
	"golang.org/x/exp/slices"
)

func ABIMethodsMapping(chain_name string, protocol string, action string, contract_address string, rpcProviderURL string) (string, string, string, string) {
	contract_mapping := map[string]map[string]map[string]string{
		"ethereum": {
			"opendao": {
				"stake":            "enter",
				"unstake":          "leave",
				"contract_address": "0xEDd27C961CE6f79afC16Fd287d934eE31a90D7D1",
				"rpcProviderURL":   "https://eth-mainnet.public.blastapi.io",
				"chain":            "ethereum",
			},
			"zeroswap": {
				"stake":            "stake",
				"unstake":          "unstake",
				"contract_address": "0xEDF822c90d62aC0557F8c4925725A2d6d6f17769",
				"rpcProviderURL":   "https://eth-mainnet.public.blastapi.io",
				"chain":            "ethereum",
			},
		},
		"polygon": {
			"zeroswap": {
				"stake":            "stake",
				"unstake":          "unstake",
				"contract_address": "0x89eA093C07f4FCc03AEBe8A1D5507c15dE88531f",
				"rpcProviderURL":   "https://polygon-rpc.com",
				"chain":            "polygon",
			},
		},
		"avax": {
			"zeroswap": {
				"stake":            "stake",
				"unstake":          "unstake",
				"contract_address": "0xa4751EAa89C5D6ff61384766268cabf25aCD1011",
				"rpcProviderURL":   "https://rpc.ankr.com/avalanche",
				"chain":            "avalanche",
			},
		},
		"bsc": {
			"zeroswap": {
				"stake":            "stake",
				"unstake":          "unstake",
				"contract_address": "0x89eA093C07f4FCc03AEBe8A1D5507c15dE88531f",
				"rpcProviderURL":   "https://bsc-dataseed.binance.org/",
				"chain":            "binance",
			},
			"pancake": {
				"stake":            "deposit",
				"unstake":          "withdraw",
				"contract_address": "0x45c54210128a065de780C4B0Df3d16664f7f859e",
				"rpcProviderURL":   "https://bsc-dataseed.binance.org/",
				"chain":            "binance",
			},
		},
	}
	return contract_mapping[chain_name][protocol][action], contract_mapping[chain_name][protocol][contract_address], contract_mapping[chain_name][protocol][rpcProviderURL], contract_mapping[chain_name][protocol]["chain"]
}

func GetLockDurationExist(protocol_name string, user_action string) bool {
	lock_duration_files := []string{"pancake"}
	exists := slices.Contains(lock_duration_files, protocol_name)
	if exists && user_action == "stake" {
		return true
	} else {
		return false
	}

}
