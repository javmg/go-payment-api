package main

import (
	"github.com/javierjmgits/go-payment-api/base/config"
)

func main() {

	configuration := config.NewConfig()

	app := NewAppStarter(configuration)

	app.Start()
}
