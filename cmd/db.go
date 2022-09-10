package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "sa123"
	dbname   = "DB_Orders"
)

var ERR_DB_NOT_FOUND = errors.New("DB: Not found")

type Database struct {
	cache *Cache
	conn  string
	db    *sql.DB
}

func CreateDatabase(cache *Cache) *Database {

	db := Database{
		conn:  "",
		cache: cache,
		db:    nil,
	}

	return &db
}

/*func (v *Database) GetOrder(id string) (Order, error) {
	/*cachedOrder := v.cache.Get(id)
	if cachedOrder != nil {
		v, ok := cachedOrder.(Order)
		if ok {
			log.Printf("Данные по заказу (id: %v) получены из КЭША", id)
			return v, nil
		}
	}

	var order Order

	var jsonObj []byte
	err := v.conn.QueryRow(context.Background(), "SELECT json FROM orders WHERE id=$1", id).Scan(&jsonObj)
	if err != nil {
		return order, err
	}

	err = json.Unmarshal(jsonObj, &order)
	if err != nil {
		return order, err
	}

	log.Printf("Данные по заказу (id: %v) получены из БАЗЫ ДАННЫХ", id)
	return order, nil
}*/

func (v *Database) RefreshCheck(ch *Cache) {
	var count int
	err := v.db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	if err != nil {
		log.Panicln(err)
	}

	if count > len(ch.data.Items) {
		for i := 1; i <= count; i++ {
			rows, err := v.db.Query("SELECT order_key, order_data FROM orders WHERE ID = $1", i)
			if err != nil {
				// handle this error better than this
				panic(err)
			}
			defer rows.Close()
			for rows.Next() {
				var order_key string
				var order_data []byte
				err = rows.Scan(&order_key, &order_data)
				if err != nil {
					// handle this error
					panic(err)
				}
				if _, found := ch.data.Items[order_key]; found {
					continue
				} else {
					var order2 Order

					err := json.Unmarshal(order_data, &order2)
					if err != nil {
						log.Panicln(err)
					}
					ch.data.Items[order_key] = order2
				}

			}
			// get any error encountered during iteration
			err = rows.Err()
			if err != nil {
				panic(err)
			}
			fmt.Println("Cache has been restored")

		}
	}
}

func (v *Database) AddOrder(order Order) {
	jsonObj, err := json.Marshal(order)

	if err != nil {
		log.Panicln(err)
	}

	insertDynStmt := `insert into "orders"("order_key", "order_data") values($1, $2)`
	_, e := v.db.Exec(insertDynStmt, order.OrderUid, jsonObj)
	if e != nil {
		log.Panicln(e)
	}

}

func (v *Database) Connect() *Database {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	v.conn = psqlconn

	if psqlconn == "" {
		log.Panicln("DB NOT FOUND")
		return nil
	}
	db, err := sql.Open("postgres", psqlconn)
	v.db = db
	if db == nil {
		log.Panicln(err)
		return nil
	}
	return v
}

func (v *Database) Close() *Database {
	defer v.db.Close()
	return v
}
