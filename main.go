package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	"wooc/config"
	"wooc/handlers"
)

func main() {
	var conf config.Config
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		log.Fatal(err)
	}

	r := router{
		ipHandler: handlers.NewIpHandler(conf),
	}

	log.Fatal(http.ListenAndServe(":8080", r.InitRouter(conf)))
}
