//-- HTTP server request/response parameters
package types

type BankMsgStruct struct {
	FromAddress string
	ToAddress string
	Amount string
}

type InstantiateMsgSturct struct {
	CodeId string
	Amount string
	Label string
	InitMsg string
	Admin string
	NoAdmin string
}

type ExecuteMsgStruct struct {
	ContractAddress string
	Amount string
	ExecMsg string
}

type QueryMsgStruct struct {
	ContractAddress string
	QueryMsg string
}

type HttpResponseStruct struct {
	ResCode int `json:"resCode"`
	ResMsg string `json:"resMsg"`
	ResData string `json:"resData"`
}