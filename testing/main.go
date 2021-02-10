package main

import (
	"fmt"
	"github.com/benka-me/yamrock/testing/config"
)

func main() {
	conf := config.NewConfig()

	fmt.Println(conf.Server().Address())
}
