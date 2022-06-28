package parse

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	cns "github.com/Moonyongjung/cns-gw/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

func ParseStoreCodeArgs(file string, sender sdk.AccAddress) (wasmtypes.MsgStoreCode, error) {
	wasm, err := ioutil.ReadFile(file)
	if err != nil {
		return wasmtypes.MsgStoreCode{}, err
	}

	// gzip the wasm file
	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)

		if err != nil {
			return wasmtypes.MsgStoreCode{}, err
		}
	} else if !ioutils.IsGzip(wasm) {
		return wasmtypes.MsgStoreCode{}, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}
	
	//-- Only sender is able to instantiate contract
	//   Terminate everybody
	var perm *wasmtypes.AccessConfig
	x := wasmtypes.AccessTypeOnlyAddress.With(sender)
	perm = &x	

	msg := wasmtypes.MsgStoreCode{
		Sender:                sender.String(),
		WASMByteCode:          wasm,
		InstantiatePermission: perm,
	}
	return msg, nil
}

func ParseInstantiateArgs(
	instantiateMsgData cns.InstantiateMsgSturct,	
	sender sdk.AccAddress) (wasmtypes.MsgInstantiateContract, error) {
	
	rawCodeID := instantiateMsgData.CodeId
	if rawCodeID == "" {
		return wasmtypes.MsgInstantiateContract{}, errors.New("No code ID")
	}
	
	// get the id of the code to instantiate
	codeID, err := strconv.ParseUint(rawCodeID, 10, 64)
	if err != nil {
		return wasmtypes.MsgInstantiateContract{}, err
	}

	amountStr := instantiateMsgData.Amount
	if amountStr == "" {
		amountStr = "0"
	}
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return wasmtypes.MsgInstantiateContract{}, fmt.Errorf("amount: %s", err)
	}

	label := instantiateMsgData.Label	
	if label == "" {
		return wasmtypes.MsgInstantiateContract{}, errors.New("label is required on all contracts")
	}

	initMsg := instantiateMsgData.InitMsg
	if initMsg == "" {
		return wasmtypes.MsgInstantiateContract{}, errors.New("No Init Message")
	}

	adminStr := instantiateMsgData.Admin

	noAdminBool := true
	noAdminStr := instantiateMsgData.NoAdmin	
	if noAdminStr == "" || noAdminStr == "true" {
		noAdminBool = true
	} else if noAdminStr == "false" {
		noAdminBool = false
	} else {
		return wasmtypes.MsgInstantiateContract{}, fmt.Errorf("noAdmin parameter must set \"true\" or \"false\"")
	}

	// ensure sensible admin is set (or explicitly immutable)
	if adminStr == "" && !noAdminBool {
		return wasmtypes.MsgInstantiateContract{}, fmt.Errorf("you must set an admin or explicitly pass --no-admin to make it immutible (wasmd issue #719)")
	}
	if adminStr != "" && noAdminBool {
		return wasmtypes.MsgInstantiateContract{}, fmt.Errorf("you set an admin and passed --no-admin, those cannot both be true")
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := wasmtypes.MsgInstantiateContract{
		Sender: sender.String(),
		CodeID: codeID,
		Label:  label,
		Funds:  amount,
		Msg:    []byte(initMsg),
		Admin:  adminStr,
	}
	return msg, nil
}

func ParseExecuteArgs(executeMsgData cns.ExecuteMsgStruct, 
	sender sdk.AccAddress) (wasmtypes.MsgExecuteContract, error){
	amountStr := executeMsgData.Amount
	if amountStr == "" {
		amountStr = "0"
	}
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return wasmtypes.MsgExecuteContract{}, fmt.Errorf("amount: %s", err)
	}

	return wasmtypes.MsgExecuteContract{
		Sender:   sender.String(),
		Contract: executeMsgData.ContractAddress,
		Funds:    amount,
		Msg:      []byte(executeMsgData.ExecMsg),
	}, nil
}

func ParseQueryArgs(queryMsgData cns.QueryMsgStruct, 
	sender sdk.AccAddress) (wasmtypes.QuerySmartContractStateRequest, error) {
		decoder := newArgDecoder(asciiDecodeString)

		queryData, err := decoder.DecodeString(queryMsgData.QueryMsg)
		if err != nil {
			return wasmtypes.QuerySmartContractStateRequest{}, errors.New(err.Error())
		}

		return wasmtypes.QuerySmartContractStateRequest {
			Address: queryMsgData.ContractAddress,
			QueryData: queryData,
		}, nil
	}