package entity

type Balance struct {
	ID       string
	UserID   string
	Balance  int
	Currency string
}

type BalanceHistory struct {
	ID                      string
	TransactionID           string // nullable
	UserID                  string
	Balance                 int
	Currency                string
	ProofImageURL           string
	CreatedAt               int64
	SourceBankAccountNumber string
	SourceBankName          string
}

type UpsertBalance struct {
	SenderBankAccountNumber string
	SenderBankName          string
	Balance                 int
	Currency                string
	ProofImageURL           string
	UserID                  string
	TransactionID           string
}

type GetBalancesHistory struct {
	Limit  int
	Offset int
	UserID string
}
