package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

type NATS struct {
	sc        stan.Conn
	NATSURL   string
	ClusterID string
	ClientID  string
}

func CreateNATS() *NATS {
	n := NATS{}

	return &n
}
func (v *NATS) Publish(obj Order) error {
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = v.sc.Publish("order", jsonObj)

	return err
}

func (v *NATS) Connect() error {
	sc, err := stan.Connect(viper.GetString("ns.ClusterID"), viper.GetString("ns.ClientID"), stan.NatsURL(viper.GetString("ns.NatsURL")))
	if err != nil {
		return err
	}
	v.sc = sc
	return nil
}

func (v *NATS) Close() {
	if v.sc != nil {
		v.sc.Close()
	}
}
