package gw

import (
	cns "github.com/Moonyongjung/cns-gw/types"
)

var responseType map[int]string

func HttpServerInit(channel cns.ChannelStruct) {
	responseType = httpResponseInit()
	go RunHttpServer(channel)
	for {
		select {
			case txRes := <- channel.TxResponseChan:
				res := httpResponseByte(int(txRes.Code), txRes.RawLog, responseType)
				channel.HttpServerChan <- res

			case queryRes := <- channel.QueryResponseChan:				
				res := httpResponseByte(0, queryRes, responseType)
				channel.HttpServerChan <- res
			
			case err := <- channel.ErrorChan:
				res := httpResponseByte(err.ResCode, err.ResData, responseType)
				channel.HttpServerChan <- res
		}
	}
}