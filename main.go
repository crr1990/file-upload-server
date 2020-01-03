package main

import (
	"file-upload-server/Router"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func main() {
	viper.SetConfigName("conf")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file error: %s\n", err)
		os.Exit(1)
	}

	Router.InitRouter()
}
