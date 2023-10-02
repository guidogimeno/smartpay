package main

import (
	"log"

	"github.com/guidogimeno/smartpay/pkg/http/rest"
)

func main() {
	app := rest.Handler()
	log.Fatal(app.Listen(":8080"))
}
