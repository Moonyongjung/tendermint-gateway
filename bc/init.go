package bc

import (
	"github.com/Moonyongjung/cns-gw/client"
	"github.com/Moonyongjung/cns-gw/bc/rest"
	"github.com/Moonyongjung/cns-gw/key"
	"github.com/Moonyongjung/cns-gw/util"
	cns "github.com/Moonyongjung/cns-gw/types"
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
				response, err := txCreate(priv, accNum, gwAdd, clientCtx, bankMsgSend, "bank")
				if err.ResCode != 0 {
					channel.ErrorChan <- err
				} else {
					channel.TxResponseChan <- response
				}

			case storeMsgSend := <- channel.StoreMsgSendChan:
				response, err := txCreate(priv, accNum, gwAdd, clientCtx, storeMsgSend, "store")
				if err.ResCode != 0 {
					channel.ErrorChan <- err
				} else {
					channel.TxResponseChan <- response
				}

			case instantiateMsgSend := <- channel.InstantiateMsgSendChan:
				response, err := txCreate(priv, accNum, gwAdd, clientCtx, instantiateMsgSend, "instantiate")
				if err.ResCode != 0 {
					channel.ErrorChan <- err
				} else {
					channel.TxResponseChan <- response
				}

			case executeMsgSend := <- channel.ExecuteMsgSendChan:
				response, err := txCreate(priv, accNum, gwAdd, clientCtx, executeMsgSend, "execute")
				if err.ResCode != 0 {
					channel.ErrorChan <- err
				} else {
					channel.TxResponseChan <- response
				}
			
			case queryMsgSend := <- channel.QueryMsgSendChan:
				response, err := querySend(clientCtx, queryMsgSend)
				if err.ResCode != 0 {
					channel.ErrorChan <- err
				} else {
					channel.QueryResponseChan <- response
				}
			}
		}
	}
}