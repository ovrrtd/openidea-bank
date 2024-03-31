package response

type Balance struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
}

type GetBalancesHistory struct {
	TransactionID    string `json:"transactionId"`
	Balance          int    `json:"balance"`
	Currency         string `json:"currency"`
	TransferProofImg string `json:"transferProofImg"`
	CreatedAt        int64  `json:"createdAt"`
	Source           struct {
		BankAccountNumber string `json:"bankAccountNumber"`
		BankName          string `json:"bankName"`
	} `json:"source"`
}
