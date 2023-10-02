package rest

import (
	"html/template"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/guidogimeno/smartpay/pkg/types"
)

type HandlerFunc func(*fiber.Ctx) error

func Handler() *fiber.App {
	app := fiber.New()

	app.Static("/views", "./pkg/views")

	app.Get("/ping", pingHandler)
	app.Get("/", smartpayHandler)
	app.Post("/analysis", analysisHandler())

	return app
}

func pingHandler(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func analysisHandler() HandlerFunc {
	type Payment struct {
		Amount       float64 `json:"amount"`
		Installments int     `json:"installments"`
		InterestRate float32 `json:"interestRate"`
	}

	return func(c *fiber.Ctx) error {
		var paymentRequest Payment
		err := c.BodyParser(&paymentRequest)
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		p := types.NewPayment(
			paymentRequest.Amount,
			paymentRequest.Installments,
			paymentRequest.InterestRate,
		)

		err = p.IsValid()
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		a, err := p.Analysis()
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		t, err := template.ParseGlob("pkg/views/analysis.html")
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		c.Type("html")
		return t.Execute(c, a)
	}
}

func smartpayHandler(c *fiber.Ctx) error {
	t, err := template.ParseGlob("pkg/views/smartpay.html")
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	c.Type("html")
	return t.Execute(c, nil)
}
