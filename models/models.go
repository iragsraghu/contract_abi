package models

// input data from user request
type InputFields struct {
	FromAddress string  `json:"from_address"`
	Amount      float64 `json:"amount"`
	Duration    float64 `json:"duration"`
	Action      string  `json:"action"`
	Chain       string  `json:"chain"`
	Protocol    string  `json:"protocol"`
}

// output data from contract
type OutputFields struct {
	EndcodedData string `json:"encoded_data"`
	Gas          int    `json:"gas"`
	From         string `json:"from"`
	To           string `json:"to"`
}
