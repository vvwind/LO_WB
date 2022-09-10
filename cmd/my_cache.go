package main

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
)

type toJson struct {
}

type Cache struct {
	mutex        sync.RWMutex
	fileName     string
	data         cacheData
	cacheCreated bool
}

type cacheData struct {
	Items  map[string]Order `json:"items"`
	myList []Order          `json:"myList"`
}

// создание пустого кеша и файла с кешем json
func (v *Cache) CreateCache(fileName string) *Cache {
	data := cacheData{
		Items: make(map[string]Order),
	}

	cache := Cache{
		mutex:        sync.RWMutex{},
		fileName:     fileName,
		data:         data,
		cacheCreated: true,
	}
	v.data.Items = make(map[string]Order)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Cant create file %s, with error %w", file, err)
	}
	return &cache
}

// добавление в кеш
func (v *Cache) addCache(key string, order Order) *Cache {

	v.mutex.RLock()
	defer v.mutex.RUnlock()
	v.data.Items[key] = order
	f, err := os.OpenFile(viper.GetString("fileName"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	b, err2 := json.Marshal(order)
	if err2 != nil {
		log.Fatal(err2)
	}

	if _, err := f.WriteString(string(b)); err != nil {
		log.Fatal(err)
	}
	f.WriteString("\n")
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return v
}

// получить из кеша
func (v *Cache) Get(key string) interface{} {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	if _, found := v.data.Items[key]; found {

		return v.data.Items[key]
	} else {
		return nil
	}

}
