package bc

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	
	"github.com/Moonyongjung/cns-gw/util"
	cns "github.com/Moonyongjung/cns-gw/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func txCreate(
	priv cryptotypes.PrivKey, 
	accNum string, 
	gwAdd sdk.AccAddress, 
	clientCtx cmclient.Context, 
	msg interface{}, 
	option string) (*sdk.TxResponse, cns.HttpResponseStruct) {

	var response cns.HttpResponseStruct

	denom := util.GetConfig().Get("denom")
	gasLimit := util.GetConfig().Get("gasLimit")
	feeAmountStr := util.GetConfig().Get("feeAmount")
	accSeq := AccSequenceMng().NowAccSeq()

	builder := clientCtx.TxConfig.NewTxBuilder()

	//-- Setting message type after transaction type check 
	if option == "bank" {
		convertMsg, _ := msg.(*bank.MsgSend)		
		builder.SetMsgs(convertMsg)

	} else if option == "store" {
		convertMsg, _ := msg.(wasm.MsgStoreCode)		
		builder.SetMsgs(&convertMsg)

	} else if option == "instantiate" {
		convertMsg, _ := msg.(wasm.MsgInstantiateContract)
		builder.SetMsgs(&convertMsg)
		
	} else if option == "execute" {
		convertMsg, _ := msg.(wasm.MsgExecuteContract)
		builder.SetMsgs(&convertMsg)
	}
	
	gasLimitu64 := util.FromStringToUint64(gasLimit)
	builder.SetGasLimit(gasLimitu64)

	feeAmounti64 := util.FromStringToInt64(feeAmountStr)

	//-- feeAmount is almost gasLimit/1000
	feeAmount := sdk.Coin{
		Amount: sdk.NewInt(feeAmounti64),
		Denom: denom,
	}
	feeAmountCoins := sdk.NewCoins(feeAmount)
	
	builder.SetFeeAmount(feeAmountCoins)	
	
	//-- If using multisig, input privs, accNums and accSeqs of other accounts
	privs := []cryptotypes.PrivKey{priv}
	accNums := []uint64{util.FromStringToUint64(accNum)}
	accSeqs := []uint64{util.FromStringToUint64(accSeq)}

	var sigsV2 []signing.SignatureV2

	sigsV2, err := txSignRound(sigsV2, clientCtx, privs, accSeqs, accNums, builder)
	if err != nil {
		util.LogGw(err)
		response.ResCode = 102
		response.ResData = err.Error()
		return nil, response
	}

	sdkTx := builder.GetTx()
	txBytes, err := clientCtx.TxConfig.TxEncoder()(sdkTx)
	if err != nil {
		util.LogGw(err)
		response.ResCode = 103
		response.ResData = err.Error()
		return nil, response
	}

	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		util.LogGw(err)
		response.ResCode = 104
		response.ResData = err.Error()
		return nil, response
	}
	util.LogGw("Transaction response", res)	

	AccSequenceMng().AddAccSeq()
	response.ResCode = 0

	return res, response
}

func txSignRound(sigsV2 []signing.SignatureV2, 
	clientCtx cmclient.Context, 
	privs []cryptotypes.PrivKey, 
	accSeqs []uint64,
	accNums []uint64,
	builder cmclient.TxBuilder) ([]signing.SignatureV2, error) {
	
	chainId := util.GetConfig().Get("chainId")

	for i, priv := range privs {
		sigV2 := signing.SignatureV2 {
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode: clientCtx.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}
		sigsV2 = append(sigsV2, sigV2)
	}

	err := builder.SetSignatures(sigsV2...)
	if err != nil {
		util.LogGw(err)
		return nil, err
	}

	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID: chainId,
			AccountNumber: accNums[i],
			Sequence: accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(
			clientCtx.TxConfig.SignModeHandler().DefaultMode(),
			signerData,
			builder,
			priv,
			clientCtx.TxConfig,
			accSeqs[i],
		)
		if err != nil {
			util.LogGw(err)
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}

	err = builder.SetSignatures(sigsV2...)
	if err != nil {
		util.LogGw(err)
		return nil, err
	}

	return sigsV2, nil
}
