package key

import (
	"io"
	"io/ioutil"
	"strings"
	"os"
	"os/exec"

	"github.com/Moonyongjung/cns-gw/util"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
)

type EncryptedJson struct {
	PriKey string 
	PubKey string 
	GwKeyAddress types.AccAddress 
}

var priKeyFileName = "pri.info"
var pubKeyFileName = "pub.info"

func NewKey() {
	keyStorePath := util.GetConfig().Get("keyStorePath")
	keyOwnerName := util.GetConfig().Get("keyOwnerName")
	keyOwnerPw := util.GetConfig().Get("keyOwnerPw")
	daemonName := util.GetConfig().Get("daemonName")

	if _, err := os.Stat(keyStorePath+"keyring-file/"); os.IsNotExist(err) {
		storeKeyringBackendFile(daemonName, keyOwnerName, keyStorePath, keyOwnerPw)		
		exportPrivKeyArmor(daemonName, keyOwnerName, keyStorePath, keyOwnerPw)

	} else {			
		util.LogGw("Key file directory :", keyStorePath)
		priKeyBytes, err := ioutil.ReadFile(keyStorePath+priKeyFileName)
		if err != nil {
			util.LogGw(err)
		}

		pubKeyBytes, err := ioutil.ReadFile(keyStorePath+pubKeyFileName)
		if err != nil {
			util.LogGw(err)
		}

		GwKey().Set(string(priKeyBytes), string(pubKeyBytes))
	}
}

func exportPrivKeyArmor(daemonName string, keyOwnerName string, keyStorePath string, keyOwnerPw string) {	
	kr, err := keyring.New("gw", keyring.BackendFile, keyStorePath, strings.NewReader(""))
	if err != nil {
		util.LogGw(err)
	}	

	pri, _ := kr.ExportPrivKeyArmor(keyOwnerName, keyOwnerPw)
	pub, _ := kr.ExportPubKeyArmor(keyOwnerName)		

	err = ioutil.WriteFile(keyStorePath+priKeyFileName, []byte(pri), 0660)
	if err != nil {
		util.LogGw(err)			
	}

	err = ioutil.WriteFile(keyStorePath+pubKeyFileName, []byte(pub), 0660)
	if err != nil {
		util.LogGw(err)			
	}

	GwKey().Set(pri, pub)
}

//-- Use daemon of each blockchain for setting prefix
func storeKeyringBackendFile(daemonName string, keyOwnerName string, keyStorePath string, keyOwnerPw string) {
	cmd := exec.Command(daemonName, "keys", "add", keyOwnerName,
			"--keyring-backend", "file",
			"--keyring-dir", keyStorePath)

	cmdExecute(cmd, keyOwnerPw)	
}

func cmdExecute(cmd *exec.Cmd, keyOwnerPw string)  {
	stdIn, _ := cmd.StdinPipe()
	// stdOut, _ := cmd.StdoutPipe()
	defer stdIn.Close()		
	
	cmd.Start()
	io.WriteString(stdIn, keyOwnerPw + "\n")
	io.WriteString(stdIn, keyOwnerPw + "\n")
	
	cmd.Wait()

	if err := cmd.Process.Kill(); err != nil {
		util.LogGw("kill", err)
	}	
}

func DeleteKeyStore() {
	keyStorePath := util.GetConfig().Get("keyStorePath")

	if _, err := os.Stat(keyStorePath); !os.IsNotExist(err) {
		os.RemoveAll(keyStorePath+"keyring-file/")
	}
}

