package client

import (
	"github.com/Moonyongjung/cns-gw/util"
	cns "github.com/Moonyongjung/cns-gw/types"
	
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/CosmWasm/wasmd/app"
)

func SetClient() cmclient.Context {
	keyStorePath := util.GetConfig().Get("keyStorePath")
	rpcEndpoint := util.GetConfig().Get("rpcEndpoint")
	chainId := util.GetConfig().Get("chainId")
	encodingConfig := cns.MakeEncodingConfig()
		
	clientCtx := cmclient.Context{}

	//-- for using resolve wasm api, need to wasm's txconfig
	clientCtx = clientCtx.WithTxConfig(app.MakeEncodingConfig().TxConfig)
	clientCtx = clientCtx.WithChainID(chainId)
	clientCtx = clientCtx.WithInterfaceRegistry(encodingConfig.InterfaceRegistry)
	clientCtx = clientCtx.WithNodeURI(rpcEndpoint)

	client, _ := cmclient.NewClientFromNode(rpcEndpoint)
	clientCtx = clientCtx.WithClient(client)

	//-- After InitKey
	clientCtx = clientCtx.WithKeyringDir(keyStorePath)
	kr, err := cmclient.NewKeyringFromBackend(clientCtx, "file")
	if err != nil {
		util.LogGw(err)
	}
	clientCtx = clientCtx.WithKeyring(kr)

	//-- To check code ID of contract, broadcast mode = block
	clientCtx = clientCtx.WithBroadcastMode("block")	
	
	return clientCtx
}
