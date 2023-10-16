package scrapper

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/guidogimeno/smartpay/pkg/http/client"
)

const (
	bcraUrl  = "https://www.bcra.gob.ar/PublicacionesEstadisticas/Principales_variables_datos.asp"
	tnaSerie = "7935"

	layout = "2006-01-02"

	tdElement = "td"
)

type BCRA struct{}

func (b *BCRA) Rate(ctx context.Context, date time.Time) (*Rate, error) {
	threeMonthsAgo := date.AddDate(0, -3, 0)
	startDate := threeMonthsAgo.Format(layout)
	finishDate := date.Format(layout)
	formData := url.Values{
		"fecha_desde": {startDate},
		"fecha_hasta": {finishDate},
		"primeravez":  {"1"},
		"serie":       {tnaSerie},
	}
	url := fmt.Sprintf("%s?serie%s", bcraUrl, tnaSerie)
	formBody := client.FormBody(formData)
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	response, err := client.Post(
		url,
		formBody,
		client.WithHeaders(headers),
		client.WithContext(ctx),
	)

	if err != nil {
		return nil, errors.New("Failed to fetch BCRA data")
	}

	index, err := parseIndex(response)
	if err != nil {
		return nil, err
	}

	return &Rate{
		Source: "BCRA",
		Index:  index,
	}, nil
}

func parseIndex(r *client.Response) (float64, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r.String()))
	if err != nil {
		return 0, err
	}

	// Get the last one
	td := doc.Find(tdElement).Eq(-1).Text()
	td = strings.TrimSpace(td)
	td = strings.ReplaceAll(td, ",", ".")

	return strconv.ParseFloat(td, 64)
}
