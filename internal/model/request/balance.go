package request

type AddBalance struct {
	BankAccountNumber string `json:"senderBankAccountNumber" validate:"required,min=5,max=30"`
	BankName          string `json:"senderBankName" validate:"required,min=3,max=30"`
	Balance           int    `json:"addedBalance" validate:"required,min=0"`
	Currency          string `json:"currency" validate:"required,iso4217"`
	ProofImageURL     string `json:"transferProofImg" validate:"required,url"`
	UserID            string
}

type GetBalancesHistory struct {
	Limit  int
	Offset int
	UserID string
}

type CreateTransaction struct {
	BankAccountNumber string `json:"recipientBankAccountNumber" validate:"required,min=5,max=30"`
	BankName          string `json:"recipientBankName" validate:"required,min=3,max=30"`
	Currency          string `json:"fromCurrency" validate:"required,iso4217"`
	Balance           int    `json:"balances" validate:"required"`
	UserID            string
}
