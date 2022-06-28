package rest

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Moonyongjung/cns-gw/util"
)

const userInfoUrl = "/cosmos/auth/v1beta1/accounts/"

//-- Get account number and sequence
func GetAccountInfoHttpClient(gwKeyAddress string) (string, string, error) {
	restEndpoint := util.GetConfig().Get("restEndpoint")
	
	url := restEndpoint + userInfoUrl + gwKeyAddress

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		util.LogHttpClient(err)
	}

	hClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	hClient.Timeout = time.Second * 30
	defer func() {
		if err := recover(); err != nil {
			util.LogHttpClient(err)		
		}
	}()
	response, err := hClient.Do(request)
	if err != nil {
		util.LogHttpClient(err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		util.LogHttpClient(err)
	}

	response.Body.Close()

	//-- The account does not have any coins or tokens, not included the chain
	if strings.Contains(string(responseBody), "code") {
		var errStruct ErrStruct
		errCheckData := util.JsonUnmarshalData(errStruct, responseBody)

		code := errCheckData.
		(map[string]interface{})["code"].(float64)

		return "", "", errors.New("Code : "+ util.ToString(code, ""))
	} else {
		var responseStruct ResponseStruct
		responseData := util.JsonUnmarshalData(responseStruct, responseBody)

		accountNumber := responseData.
		(map[string]interface{})["account"].
		(map[string]interface{})["account_number"].(string)

		sequence := responseData.
		(map[string]interface{})["account"].
		(map[string]interface{})["sequence"].(string)

		return accountNumber, sequence, nil		
	}
}
