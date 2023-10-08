package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/guidogimeno/smartpay/pkg/types"
)

type HandlerFunc func(*fiber.Ctx) error

func Handler() *fiber.App {
	engine := html.New("./public/templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/public", "./public")

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
			return c.Render("error", "Invalid payment")
		}

		p := types.NewPayment(
			paymentRequest.Amount,
			paymentRequest.Installments,
			paymentRequest.InterestRate,
		)

		err = p.IsValid()
		if err != nil {
			return c.Render("error", err.Error())
		}

		a, err := p.Analysis()
		if err != nil {
			return c.Render("error", err.Error())
		}

		return c.Render("analysis", a)
	}
}

func smartpayHandler(c *fiber.Ctx) error {
	return c.Render("smartpay", nil)
}
