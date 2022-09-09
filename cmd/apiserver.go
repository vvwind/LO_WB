package main

import (
	"bufio"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"os"
)

type API struct {
	nats *NATS
}

func CreateAPISERVER(nats *NATS) *API {
	api := API{
		nats: nats,
	}

	return &api
}

func (a *API) Run() error {

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.StaticFile("/home", "./front/index.html")
	router.POST("/home/", a.OrderHandler)
	err := router.Run(viper.GetString("port"))

	return err
}

func (a *API) OrderHandler(c *gin.Context) {
	for {
		reader := bufio.NewReader(os.Stdin)
		var jsonData []byte
		for {
			jsonD, err := reader.ReadString('\n')
			if len(jsonD) == 2 {
				break
			}
			if err != nil {
				log.Println(err)
				break
			}
			b := []byte(jsonD)
			jsonData = append(jsonData, b...)
		}
		var order Order
		err := json.Unmarshal(jsonData, &order)
		if err != nil {
			log.Println(err)
			break
		}
		err2 := a.nats.Publish(order)
		if err2 != nil {
			log.Println(err2)
			break
		}
		log.Println("Успешно отправлены данные")
	}

}
