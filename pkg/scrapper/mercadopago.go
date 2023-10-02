package scrapper

import (
	"errors"
	"regexp"
	"smartpay/pkg/http/client"
	"strconv"
	"time"
)

const (
	urlMercadoPago = "https://www.mercadopago.com.ar/cuenta"
)

type MercadoPago struct{}

func (r *MercadoPago) Rate(date time.Time) (*Rate, error) {
	response, err := client.Get(urlMercadoPago)
	if err != nil {
		return nil, err
	}

	index, err := r.parseIndex(response.String())
	if err != nil {
		return nil, err
	}

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
