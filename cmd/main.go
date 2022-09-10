package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("err initing config %s", err.Error())
	}

	// cache
	cache := Cache{}
	cache.CreateCache(viper.GetString("fileName"))

	// nats
	nats := CreateNATS()
	err0 := nats.Connect()
	if err0 != nil {
		log.Println(err0)
	}
	defer nats.Close()

	// db
	db := CreateDatabase(&cache)
	db = db.Connect()
	if db == nil {
		log.Panicln(db)
	}

	// data
	err2 := nats.Subscribe(&cache, db)
	time.Sleep(7 * time.Second)
	if err2 != nil {
		fmt.Printf("%s", err2)
	}
	db.RefreshCheck(&cache)

	defer db.Close()

	api := CreateAPI(db, &cache)
	errapi := api.Run()
	if errapi != nil {
		log.Panicln(errapi)
	}
	fmt.Scanln() // waiting

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
