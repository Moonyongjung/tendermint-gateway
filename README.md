# Gateway based on Cosmos SDK
HTTP API gateway for blockchain based on Tendermint core. It can support basic functionalities of Smart contract (by CosmWasm) 

## Prerequisites
- Download blockchain daemon(e.g. gaiad, terrad, etc) or self-build tendermint based blockchain by using Ignite CLI in order to generate gateway address for matching prefix. 
- Set config file
  - `./config/config.json`
  - 
  ```yaml
  {
    "gatewayServerPort": "12000",
    "chainId": "localnet",
    "restEndpoint": "http://localhost:1317",
    "rpcEndpoint": "http://localhost:26657",
    "daemonName": "nonamed",
    "denom": "unoname",
    "contractPath": "./wasm/",
    "contractName": "contract.wasm",
    "bech32Prefix": "noname",
    "gasLimit": "5000000",
    "feeAmount": "5000"
  }
  ```
  - `gatewayServerPort` : Gateway server TCP port number
  - `chainId`, `restEndpoint`, `rpcEndpoint` : Target chain information
  - `daemonName` : Downloaded or self-builded daemon name
  - `denom` : Currency denomination (e.g. uatom)
  - `contractPath`, `contractName` : Contract local path and name to deploy
  - `bech32Prefix` : Prefix of account address and smart contract address
  - `gasLimit`, `feeAmount` : Fees to send transactions
- Set key config file
  - `./config/configKey.json`
  - `keyStorePath` : Keyring backend path (--keyring-backend file, default path: ./key/)
  - `keyOwnerName`, `keyOwnerPw` : Gateway key owner. `keyOwnerPw` must be the same as the password entered in the cli when generating the initial key. (i.e. Enter keyring passpharse of Cosmos SDK)
  
## Start
```shell
go mod tidy
go build gw.go
./gw
```

## API
All Parameters of request json value are string type
### Coin send (test)
  - (POST) `/api/bank/send`
  - Request body
  ``` yaml
  {
        "fromAddress": "a",
        "toAddress": "b",
        "amount": "100"
  }
  ```
  - `fromAddress` : GW address
### Store contract
  - (GET) `/api/wasm/store`
  - Request after copy .wasm file to the `contractPath`
### Instantiate 
  - (POST) `/api/wasm/instantiate`
  - Request body
  ```yaml
  {
        "codeId": "1",
        "amount": "1000unoname",
        "label": "contract inst",
        "initMsg": "{\"purchase_price\":{\"amount\":\"100\",\"denom\":\"unoname\"}"
  }
  ```
  - `codeId` : Contract code ID
### Execute
  - (POST) `/api/wasm/execute`
  - Request body
  ```yaml
  {
        "contractAddress": "noname19h0d6k4mtxw5qjr0aretjy9kwyem0hxclf88ka2uwjn47e90mqrqk4tkjt",
        "amount": "0unoname",
        "execMsg": "{\"register\":{\"name\":\"fred\"}}"
  }
  ```
### Query - Contract state
  - (POST) `/api/wasm/query`
  - Request body
  ```yaml
  {
        "contractAddress": "noname19h0d6k4mtxw5qjr0aretjy9kwyem0hxclf88ka2uwjn47e90mqrqk4tkjt",
        "queryMsg": "{\"resolve_record\": {\"name\": \"fred\"}}"
  }
  ```
### Query - Contract list
  - (GET) `/api/wasm/list-code`
### Query - Contract information by code ID
  - (POST) `/api/wasm/list-contract-by-code`
  - Request body
  ```yaml
  {
        "codeId": "1"
  }
  ```
### Query - Download contract wasm file
  - (POST) `/api/wasm/download`
  - Request body
  ```yaml
  {
        "codeId": "1",
        "downloadFileName":"download"
  }
  ```
### Query - Code information for a given code ID
  - (POST) `/api/wasm/code-info`
  - Request body
  ```yaml
  {
        "codeId": "1"
  }
  ```
### Query - Contract information for a given contract address
  - (POST) `/api/wasm/contract-info`
  - Request body
  ```yaml
  {
        "contractAddress": "noname19h0d6k4mtxw5qjr0aretjy9kwyem0hxclf88ka2uwjn47e90mqrqk4tkjt"
  }
  ```
### Query - All of contract internal state
  - (POST) `/api/wasm/contract-state-all`
  - Request body
  ```yaml
  {
        "contractAddress": "noname19h0d6k4mtxw5qjr0aretjy9kwyem0hxclf88ka2uwjn47e90mqrqk4tkjt"
  }
  ```
### Query - Contract history
  - (POST) `/api/wasm/contract-history`
  - Request body
  ```yaml
  {
        "contractAddress": "noname19h0d6k4mtxw5qjr0aretjy9kwyem0hxclf88ka2uwjn47e90mqrqk4tkjt"
  }
  ```
### Query - Pinned code
  - (Get) `/api/wasm/pinned`  
 ### Gatway HTTP Response 
  - `resCode` is int type, `resMsg` and `resData` is string type
  - Standard
  ```yaml
  {
        "resCode": 0,
        "resMsg": "",
        "resData": ""
  }
  ```
