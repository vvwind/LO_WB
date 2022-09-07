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

	router.StaticFile("/", "./front/index.html")

	//router.POST("/api/order", a.postOrderHandler)
	router.GET("/orderf", a.postOrderHandler)
	err := router.Run(viper.GetString("port"))

	return err
}
func (v *API) postOrderHandler(c *gin.Context) {
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
		return
	}
	err2 := v.nats.Publish(order)
	if err2 != nil {
		log.Println(err2)
		return
	}
	log.Println("Успешно отправлены данные")
}

/*func (v *API) postOrderHandler(c *gin.Context) {
	jsonData, err := c.GetRawData()
	fmt.Print(jsonData)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"message": "looser"})
		return
	}

	var order Order
	err = json.Unmarshal(jsonData, &order)
	if err != nil {
		c.JSON(400, gin.H{"message": "looser"})
		return
	}

	err = v.nats.Publish(order)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{})
		return
	}

	c.JSON(201, gin.H{})
}*/
