package main

import (	
	"os"
	"os/signal"
	"syscall"

	"github.com/Moonyongjung/cns-gw/bc"
	"github.com/Moonyongjung/cns-gw/util"
	"github.com/Moonyongjung/cns-gw/gw"
	"github.com/Moonyongjung/cns-gw/key"
	cns "github.com/Moonyongjung/cns-gw/types"
)

var configPath = "./config/config.json"
var configKeyPath = "./config/configKey.json"

var Channel cns.ChannelStruct

func init() {
	util.GetConfig().Read(configPath)
	util.GetConfig().Read(configKeyPath)
	util.SetChainPrefixConfig()	
	key.NewKey()
}

func main() {
	Channel := cns.ChannelInit()
	
	go bc.TxInit(Channel)
	go gw.RunHttpServer(Channel)
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	util.LogGw("Shutting down the server...")
	util.LogGw("Server gracefully stopped")	
}