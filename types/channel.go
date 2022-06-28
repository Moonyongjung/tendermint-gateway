package types

import (
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
)

type ChannelStruct struct {
	BankMsgSendChan chan *bank.MsgSend
	StoreMsgSendChan chan wasm.MsgStoreCode
	InstantiateMsgSendChan chan wasm.MsgInstantiateContract
	ExecuteMsgSendChan chan wasm.MsgExecuteContract
	QueryMsgSendChan chan wasm.QuerySmartContractStateRequest
	TxResponseChan chan *sdk.TxResponse
	QueryResponseChan chan wasm.QuerySmartContractStateResponse	
	HttpServerChan chan []byte
	ErrorChan chan HttpResponseStruct
}

func ChannelInit() ChannelStruct{
	var channel ChannelStruct
	channel.BankMsgSendChan = make(chan *bank.MsgSend)
	channel.StoreMsgSendChan = make(chan wasm.MsgStoreCode)
	channel.InstantiateMsgSendChan = make(chan wasm.MsgInstantiateContract)
	channel.ExecuteMsgSendChan = make(chan wasm.MsgExecuteContract)
	channel.QueryMsgSendChan = make(chan wasm.QuerySmartContractStateRequest)
	channel.TxResponseChan = make(chan *sdk.TxResponse)
	channel.QueryResponseChan = make(chan wasm.QuerySmartContractStateResponse)
	channel.HttpServerChan = make(chan []byte)
	channel.ErrorChan = make(chan HttpResponseStruct)

	return channel
}