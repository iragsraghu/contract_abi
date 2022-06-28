package models

type InputFields struct {
	FromAddress string  `json:"from_address"`
	Amount      float64 `json:"amount"`
	Duration    float64 `json:"duration"`
	Action      string  `json:"action"`
	Chain       string  `json:"chain"`
	Protocol    string  `json:"protocol"`
}
