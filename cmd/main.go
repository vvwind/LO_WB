package main

import (
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("err initing config %s", err.Error())
	}

	nats := CreateNATS()
	err := nats.Connect()
	if err != nil {
		log.Panicln(err)
	}
	defer nats.Close()

	api := CreateAPISERVER(nats)
	err2 := api.Run()
	if err2 != nil {
		log.Panicln(err2)
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
