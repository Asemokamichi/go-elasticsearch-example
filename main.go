package main

import (
	"context"
	"io/ioutil"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/elastic/go-elasticsearch/esapi"
)

// предварительно надо поменять данные
const (
	username = "elastic"
	password = "_BGma68QLk_SwwzgJf=R"
	fileCert = "http_ca.crt"
)

var (
	addresses = []string{"https://localhost:9200"}
)

func main() {
	ctx := context.Background()
	cert, err := ioutil.ReadFile(fileCert)

	//Заполняем конфигурационные данные
	cfg := elasticsearch.Config{
		Addresses: addresses, //	Список узлов Elasticsearch для использования
		Username:  username,  //	Имя пользователя для базовой аутентификации по протоколу HTTP.
		Password:  password,  //	Пароль для базовой аутентификации по протоколу HTTP.
		CACert:    cert,      //	сертификация
	}

	//	Создаем нового клиента
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("Elasticsearch connection error:", err)
	}

	// Проверяем подключение, если все правильно отправить ответ запроса
	info, err := client.Info()
	if err != nil {
		log.Fatalln("client response error: ", err)
	} else {
		log.Println("client response: ", info)
	}

	//  Настраиваем запрос API переиндексации
	req := esapi.ReindexRequest{
		Body: strings.NewReader(`{
		"source": {
			"index": "eventlog",
			"query": {
				"range": {
					"timestamp": {
						"gte": "2023-02-01T10:00:00",
						"lte": "2023-03-01T10:00:00"
					}
				}
			}
	
		},
		"dest": {
			"index": "log"
		}
	}`),
	}

	// С помощью функции Do() отправляем запрос req, ответ записываем rr
	// проверяем на ошибку, если все в порядке то проверяем ответ res
	res, err := req.Do(ctx, client)
	if err != nil {
		log.Fatalf("IndexRequest ERROR: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalln("client response error: ", err)
	} else {
		log.Println("response: ", res)
	}

}
