package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
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

	for {
		err := OrderHandler(nats)
		if err != nil {
			log.Println(err)
		}

	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func OrderHandler(v *NATS) error {

	reader := bufio.NewReader(os.Stdin)
	var jsonData []byte
	for {
		jsonD, err := reader.ReadString('\n')
		if len(jsonD) == 2 {

			break
		}
		if err != nil {
			return err
		}
		b := []byte(jsonD)
		jsonData = append(jsonData, b...)
	}
	var order Order
	err := json.Unmarshal(jsonData, &order)
	if err != nil {
		return err
	}
	if order.OrderUid == "" {
		err1 := errors.New("Can't create order with no ID")
		return err1
	}
	err2 := v.Publish(order)
	if err2 != nil {
		return err2
	}
	log.Println("Data has been send!")
	return nil
}
