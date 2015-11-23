package conf

import (
	"github.com/Unknwon/goconfig"
	"os"
	"../helper"
)

var (
	config *goconfig.ConfigFile
)

func InitConfig() error {
	var err error
	config, err = goconfig.LoadConfigFile("./config.ini")
	if err != nil {
		return err
	}
	
	return nil
}

func CreateIniFile() error {
	f, err := os.OpenFile("./server.ini", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	defer f.Close()
	
	config, err = goconfig.LoadConfigFile("./server.ini")
	if err != nil {
		return err
	}

	config.SetSectionComments("App", "")
	config.SetSectionComments("User", "")

	config.SetValue("App", "IP", "127.0.0.1")
	config.SetValue("App", "Port", "8080")
	config.SetValue("User", "username", "admin")
	config.SetValue("User", "password", helper.MD5("admin"))

	err = goconfig.SaveConfigFile(config, "./server.ini")
	if err != nil {
		return err
	}

	return nil
}

func ReadKey(section, key string) (string, error) {
	return config.GetValue(section, key)
}

func WriteKey(section, key, value string) bool {
	return config.SetValue(section, key, value)
}