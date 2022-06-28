//-- Get account info parameters
package rest

type ErrStruct struct {
	Code float64
	Message string
	Details string
}

type ResponseStruct struct{
	Account AccountStruct
}

type AccountStruct struct {
	Type string
	Address string
	Pubkey PubKeyStruct
	AccountNumber string
	Sequence string
}

type PubKeyStruct struct {
	Type string
	Key string
}