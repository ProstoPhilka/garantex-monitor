package main

import (
	"fmt"
	"garantex-monitor/config"
)

func main() {

	cfg := config.MustLoad(".env")

	fmt.Println(cfg)
}
