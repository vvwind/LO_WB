package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
	"log"
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

func (v *NATS) Subscribe(ch *Cache, db *Database) error {
	_, err := v.sc.Subscribe("order", func(m *stan.Msg) {
		if err := m.Ack(); err != nil {
			log.Printf("failed to ACK msg:%v", err)
			return
		}
		var order Order
		fmt.Println(string(m.Data))
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			log.Println(err)
			return
		}
		if order.OrderUid != "" {
			ch = ch.addCache(order.OrderUid, order)
			db.AddOrder(order)
		}

	}, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.DurableName("my-durable"))

	return err
}

func (v *NATS) Connect() error {

	sc, err := stan.Connect(viper.GetString("ns.ClusterID"), viper.GetString("ns.ClientID"), stan.NatsURL(viper.GetString("ns.NatsURL")), stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		log.Fatalf("Connection lost, reason: %v", reason)
	}))
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
