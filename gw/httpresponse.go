package gw

import (
	cns "github.com/Moonyongjung/cns-gw/types"
	"github.com/Moonyongjung/cns-gw/util"
)

//-- Response code
//   If response code is not included in below response code, it is Cosmos SDK or Wasm code

var response map[int]string

func httpResponseInit() map[int]string {
	response = make(map[int]string)

	response[0] = "Success"
	// response[1] = "Wrong API version"
	response[102] = "Transaction signing error"
	response[103] = "Transaction encoding error"
	response[104] = "Transaction broadcast error"
	response[105] = "Query error"
	response[106] = "Message creating error"

	return response

}
func httpResponseByte(resCode int, resData string, responseType map[int]string) ([]byte) {
	
	var httpResponse cns.HttpResponseStruct

	httpResponse.ResCode = resCode
	httpResponse.ResMsg = responseType[resCode]
	httpResponse.ResData = resData

	responseByte, err := util.JsonMarshalData(httpResponse)
	if err != nil {
		util.LogGw(err)
	}

	return responseByte
}