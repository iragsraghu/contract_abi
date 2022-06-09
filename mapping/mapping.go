package mapping

func MappedMethods(abi_file_name string, user_action string) string {
	explorers := map[string]map[string]string{
		// Etherscan chain
		"opendao": {
			"stake":   "enter",
			"unstake": "leave",
		},
		// BSC chain
		"pancake": {
			"stake":   "deposit",
			"unstake": "leave",
		},
		// Polygon chain
		"beefy": {
			"stake":   "stake",
			"unstake": "withdraw",
		},
	}
	return explorers[abi_file_name][user_action]
}
