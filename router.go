package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"wooc/config"
	"wooc/handlers"
	"wooc/middlewares"
)

type router struct {
	ipHandler *handlers.IpHandler
}

func (router router) InitRouter(config config.Config) *mux.Router {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./swaggerui/"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	r.HandleFunc("/ip", router.ipHandler.GetList).Methods(http.MethodGet)
	r.HandleFunc("/ip", router.ipHandler.Add).Methods(http.MethodPost)

	amw := middlewares.AuthMiddleware{Tokens: config.Tokens}
	r.Use(amw.Middleware)

	return r
}