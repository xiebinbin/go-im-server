package config

import (
	"imsdk/internal/client/pkg/utils/filehelper"
	walletutil "imsdk/internal/client/pkg/utils/wallet-util"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

var c = AppConfig{}
var serverWallet *walletutil.Wallet

type AppConfig struct {
	PriKey string `yaml:"priKey"`
	DbFile string `yaml:"dbFile"`
	Port   string `yaml:"port"`
}

func InitConfig(file string) {

	fileInfo, err := os.Stat(file)
	if err != nil {
		panic(err)
	}
	buffer := make([]byte, fileInfo.Size())
	fileHandler, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	_, err = io.ReadFull(fileHandler, buffer)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(buffer, &c)
	if err != nil {
		panic(err)
	}
	c.DbFile = filehelper.GetAppPath(c.DbFile)

	c.PriKey = filehelper.GetAppPath(c.PriKey)

}
func InitPriKey() {
	//if fileutil.IsExist(c.PriKey) == false {
	//	key, err := walletutil.GenerateKey()
	//	if err != nil {
	//		panic(err)
	//	}
	//	err = fileutil.WriteBytesToFile(c.PriKey, []byte(key))
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	keyBytes, err := ioutil.ReadFile(c.PriKey)
	if err != nil {
		panic(err)
	}
	wallet, err := walletutil.New(string(keyBytes))
	if err != nil {
		panic(err)
	}
	serverWallet = wallet
}
