package models

type RequiredFields struct {
	ContractAddress string `json:"contract_address"`
	RPC             string `json:"rpc"`
	Action          string `json:"action"`
	Protocol        string `json:"protocol"`
	Chain           string `json:"chain"`
	WalletAddress   string `json:"wallet_address"`
}

type InputFields struct {
	FromAddress string  `json:"from_address"`
	Amount      float64 `json:"amount"`
	Duration    float64 `json:"duration"`
	Action      string  `json:"action"`
	Chain       string  `json:"chain"`
	Protocol    string  `json:"protocol"`
}
