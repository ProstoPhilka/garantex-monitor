package main

import (
	"garantex-monitor/config"
	"garantex-monitor/internal/infrastructure/logs"

	"go.uber.org/zap"
)

func main() {

	conf := config.MustLoad()

	logger := logs.NewLogger(conf)

	logger.Info("logger", zap.String("key", "value"))
}

// package garantexapi

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// // type GarantexApi struct {
// // 	url string
// // }

// func GetRequest() {
// 	resp, err := http.Get("https://garantex.org/api/v2/trades?market=usdtrub")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Чтение тела ответа
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Выводим содержимое тела ответа
// 	fmt.Println(string(body))
// }
