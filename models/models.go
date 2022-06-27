package models

type RequiredFields struct {
	ContractAddress string `json:"contract_address"`
	RPC             string `json:"rpc"`
	Action          string `json:"action"`
	Protocol        string `json:"protocol"`
	Chain           string `json:"chain"`
	WalletAddress   string `json:"wallet_address"`
}
