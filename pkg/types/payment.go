package types

import (
	"errors"
	"time"

	"github.com/guidogimeno/smartpay/pkg/scrapper"
	"github.com/shopspring/decimal"
)

const (
	months = 12
)

type Payment struct {
	Amount       float64
	Installments int
	InterestRate float32
}

func NewPayment(amount float64, installments int, interestRate float32) *Payment {
	return &Payment{
		Amount:       amount,
		Installments: installments,
		InterestRate: interestRate,
	}
}

func (p *Payment) DecimalAmount() decimal.Decimal {
	return decimal.NewFromFloat(p.Amount)
}

func (p *Payment) IsValid() error {
	if p.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if p.InterestRate < 0 {
		return errors.New("interest rate must be positive")
	}

	if p.Installments < 0 {
		return errors.New("number of installments must be positive")
	}

	return nil
}

func (p *Payment) Analysis() ([]*Analysis, error) {
	var analysis []*Analysis
	ratables := []scrapper.Ratable{
		&scrapper.MercadoPago{},
		&scrapper.BCRA{},
	}

	for _, r := range ratables {
		a, err := p.doTheMath(r)
		if err != nil {
			return nil, err
		}

		analysis = append(analysis, a)
	}

	return analysis, nil
}

func (p *Payment) doTheMath(r scrapper.Ratable) (*Analysis, error) {
	today := time.Now()
	rate, err := r.Rate(today)
	if err != nil {
		return nil, err
	}

	tna := decimal.NewFromFloat(rate.Index).Div(decimal.NewFromInt(100))
	tnm := tna.Div(decimal.NewFromInt(months))
	installmentWithInterest := p.installmentWithInterest()

	savings := p.DecimalAmount().Copy()
	for i := 1; i <= p.Installments; i++ {
		fixedDepositInterest := savings.Mul(tnm)
		savings = savings.Add(fixedDepositInterest).Sub(installmentWithInterest)
	}

	return &Analysis{
		Entity:  rate.Source,
		Savings: savings.RoundDown(0).InexactFloat64(),
		Index:   rate.Index,
	}, nil
}

func (p *Payment) installmentWithInterest() decimal.Decimal {
	installments := decimal.NewFromInt(int64(p.Installments))
	interestRate := decimal.NewFromFloat32(p.InterestRate)

	installmentAmount := p.DecimalAmount().Div(installments)
	interest := installmentAmount.Mul(interestRate)

	return installmentAmount.Add(interest)
}
