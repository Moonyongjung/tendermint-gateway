package bc

import (
	"context"

	"github.com/Moonyongjung/cns-gw/util"
	cns "github.com/Moonyongjung/cns-gw/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

)

//-- Query contract state using by NewQueryClient of wasm
func querySend(clientCtx cmclient.Context, 
	queryMsgSend wasmtypes.QuerySmartContractStateRequest) (wasmtypes.QuerySmartContractStateResponse, cns.HttpResponseStruct){
	
	var response cns.HttpResponseStruct

	queryClient := wasmtypes.NewQueryClient(clientCtx)
	res, err := queryClient.SmartContractState(
		context.Background(),
		&queryMsgSend,
	)
	if err != nil {
		response.ResCode = 105
		response.ResData = err.Error()
		return wasmtypes.QuerySmartContractStateResponse{}, response
	}

	util.LogGw("Transaction response", res)	

	response.ResCode = 0
	return *res, response
}