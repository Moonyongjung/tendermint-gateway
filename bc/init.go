package bc

import (
	"github.com/Moonyongjung/cns-gw/client"
	"github.com/Moonyongjung/cns-gw/bc/rest"
	"github.com/Moonyongjung/cns-gw/key"
	"github.com/Moonyongjung/cns-gw/util"
	cns "github.com/Moonyongjung/cns-gw/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)
func TxInit(channel cns.ChannelStruct) {	
	//-- Cosmos SDK client init
	clientCtx := client.SetClient()

	//-- Get Private key and GW address
	priv := util.GetPriKeyByArmor(key.GwKey().GetPriKey())
	gwAdd := util.GetGwAddrByPrivKey(priv)
	util.LogGw("GW address : ", gwAdd.String())
	
	//-- Get GW address's account number and sequence
	accNum, accSeq, err := rest.GetAccountInfoHttpClient(gwAdd.String())

	//-- Setting account num and seq
	util.GetConfigAcc().Set(accNum, accSeq)
	AccSequenceMng().NewAccSeq()

	//-- If gw account address does not included chain, err return
	if err != nil {
		util.LogGw(err)
		util.LogGw("account "+gwAdd.String()+" not found")
		util.LogGw("Get coin for using tx fee")
	} else {
		//-- Msg wait..
		for {
			select{
			case bankMsgSend := <- channel.BankMsgSendChan:
				txResponse, httpResponse := txCreate(priv, accNum, gwAdd, clientCtx, bankMsgSend, "bank")
				txResponseFunc(txResponse, httpResponse, channel)

			case storeMsgSend := <- channel.StoreMsgSendChan:
				txResponse, httpResponse := txCreate(priv, accNum, gwAdd, clientCtx, storeMsgSend, "store")
				txResponseFunc(txResponse, httpResponse, channel)

			case instantiateMsgSend := <- channel.InstantiateMsgSendChan:
				txResponse, httpResponse := txCreate(priv, accNum, gwAdd, clientCtx, instantiateMsgSend, "instantiate")
				txResponseFunc(txResponse, httpResponse, channel)

			case executeMsgSend := <- channel.ExecuteMsgSendChan:
				txResponse, httpResponse := txCreate(priv, accNum, gwAdd, clientCtx, executeMsgSend, "execute")
				txResponseFunc(txResponse, httpResponse, channel)
			
			case queryMsgSend := <- channel.QueryMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, queryMsgSend, "query")
				queryResponseFunc(queryResponse, httpResponse, channel)
			
			case listcodeMsgSend := <- channel.ListcodeMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, listcodeMsgSend, "listcode")
				queryResponseFunc(queryResponse, httpResponse, channel)
			
			case listContractByCodeMsgSend := <- channel.ListContractByCodeMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, listContractByCodeMsgSend, "listcontractbycode")
				queryResponseFunc(queryResponse, httpResponse, channel)
				
			case downloadMsgSend := <- channel.DownloadMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, downloadMsgSend, "download")
				queryResponseFunc(queryResponse, httpResponse, channel)
			
			case codeInfoMsgSend := <- channel.CodeInfoMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, codeInfoMsgSend, "codeinfo")
				queryResponseFunc(queryResponse, httpResponse, channel)
			
			case contractInfoMsgSend := <- channel.ContractInfoMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, contractInfoMsgSend, "contractinfo")
				queryResponseFunc(queryResponse, httpResponse, channel)

			case contractStateAllMsgSend := <- channel.ContractStateAllMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, contractStateAllMsgSend, "contractstateall")
				queryResponseFunc(queryResponse, httpResponse, channel)

			case contractHistoryMsgSend := <- channel.ContractHistoryMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, contractHistoryMsgSend, "contracthistory")
				queryResponseFunc(queryResponse, httpResponse, channel)

			case pinnedMsgSend := <- channel.PinnedMsgSendChan:
				queryResponse, httpResponse := querySend(clientCtx, pinnedMsgSend, "pinned")
				queryResponseFunc(queryResponse, httpResponse, channel)
			}
		}
	}
}

func txResponseFunc(txResponse *sdk.TxResponse, 
	httpResponse cns.HttpResponseStruct,
	channel cns.ChannelStruct) {

	if httpResponse.ResCode != 0 {
		channel.ErrorChan <- httpResponse
	} else {
		channel.TxResponseChan <- txResponse
	}
}

func queryResponseFunc(queryResponse string, 
	httpResponse cns.HttpResponseStruct,
	channel cns.ChannelStruct) {

	if httpResponse.ResCode != 0 {
		channel.ErrorChan <- httpResponse
	} else {
		channel.QueryResponseChan <- queryResponse
	}
}