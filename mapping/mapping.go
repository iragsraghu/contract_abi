package mapping

func ABIMethodsMapping(abi_file_name string, user_action string) string {
	explorers := map[string]map[string]string{
		// Etherscan chain
		"opendao": {
			"stake":   "enter",
			"unstake": "leave",
		},
		// BSC chain
		"pancake": {
			"stake":   "deposit",
			"unstake": "withdraw",
		},
		// Polygon chain
		"beefy": {
			"stake":   "stake",
			"unstake": "withdraw",
		},
	}
	return explorers[abi_file_name][user_action]
}

// getting chain name from contract address
func GetChainName(contract_address string) string {
	abi_file_names := map[string]string{
		// Etherscan chain
		"0xEDd27C961CE6f79afC16Fd287d934eE31a90D7D1": "opendao",
		// BSC chain
		"0x45c54210128a065de780C4B0Df3d16664f7f859e": "pancake",
		// Polygon chain
		"0xDeB0a777ba6f59C78c654B8c92F80238c8002DD2": "beefy",
	}
	return abi_file_names[contract_address]
}

func GetLockDurationExist(abi_file_name string) []string {
	return []string{"pancake"} // add abi filename if lock duration exist in that chain
}
