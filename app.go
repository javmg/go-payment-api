package main

import (
	"gitgub.com/javierjmgits/go-payment-api/base/config"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"log"
	"net/http"
	"gitgub.com/javierjmgits/go-payment-api/payment/handler"
	"gitgub.com/javierjmgits/go-payment-api/payment/model"
	"gitgub.com/javierjmgits/go-payment-api/payment/repository"
)

type app struct {
	config *config.Config
}

type AppStarter interface {
	Start()
}

func NewAppStarter(config *config.Config) AppStarter {

	return &app{
		config: config,
	}
}

func (app *app) Start() {

	log.Println("Starting the application...")

	//
	// DB

	dbURL := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True",
		app.config.DB.Username,
		app.config.DB.Password,
		app.config.DB.Name, )

	db, err := gorm.Open("mysql", dbURL)

	if err != nil {
		log.Fatal("Error connecting to DB", err)
	}

	log.Println("Connection established with DB")

	defer db.Close()

	db = model.SetUp(db)

	//
	// Routing

	router := mux.NewRouter()

	handler.NewPaymentHandler(repository.NewPaymentRepositoryImpl(db)).Register(router)

	//
	// Server

	address := fmt.Sprintf("%v:%v", app.config.Server.Host, app.config.Server.Port)

	log.Printf("Server listening at: %v\n", address)

	errorListenAndServe := http.ListenAndServe(address, router)

	log.Fatal(errorListenAndServe)
}
