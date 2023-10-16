package scrapper

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/guidogimeno/smartpay/pkg/http/client"
)

const (
	urlMercadoPago = "https://www.mercadopago.com.ar/cuenta"
)

type MercadoPago struct{}

func (r *MercadoPago) Rate(ctx context.Context, date time.Time) (*Rate, error) {
	response, err := client.Get(urlMercadoPago, client.WithContext(ctx))
	if err != nil {
		log.Println("error mp", err)
		return nil, err
	}

	index, err := r.parseIndex(response.String())
	if err != nil {
		return nil, err
	}
	log.Println("index mp", index)

	return &Rate{
		Source: "Mercado Pago",
		Index:  index,
	}, nil
}

func (r *MercadoPago) parseIndex(response string) (float64, error) {
	pattern := "Rinde ([0-9.]+)%"
	reg := regexp.MustCompile(pattern)

	matches := reg.FindAllString(response, -1)
	if len(matches) == 0 {
		return 0, errors.New("MercadoPago rate not found")
	}

	pattern = `[0-9]+(?:\.[0-9]+)?`
	reg = regexp.MustCompile(pattern)

	match := reg.FindString(matches[0])

	return strconv.ParseFloat(match, 64)
}
