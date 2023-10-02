package main

import (
	"log"
	"smartpay/pkg/http/rest"
)

func main() {
	app := rest.Handler()
	log.Fatal(app.Listen(":8080"))
}
