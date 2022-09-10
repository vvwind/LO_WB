package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type API struct {
	db *Database
	ch *Cache
}

func CreateAPI(db *Database, ch *Cache) *API {
	api := API{
		ch: ch,
		db: db,
	}

	return &api
}

func (v *API) Run() error {

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.StaticFile("/", "./front/index.html")

	router.GET("/api/order/:id", v.getOrderHandler)

	err := router.Run(viper.GetString("port"))

	return err
}

func (v *API) getOrderHandler(c *gin.Context) {
	orderid := v.ch.Get(c.Param("id"))
	if orderid == nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, orderid)
}
