package config

import (
	"encoding/json"
	"os"
	"fmt"
)

type Configuration struct {
	DBPORT     		int
	APP_PORT        string
	LSPATH    	    string
	LSURL           string
	IP_CHECK		bool
	QIWI_IP_1       string
	QIWI_IP_2		string
	LS_LOGIN		string
	LS_PASSWORD		string
}

var Config *Configuration

func LoadConfig() {
	file, err := os.Open("./config/conf.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = &Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error decoding config:", err)
	}
}
