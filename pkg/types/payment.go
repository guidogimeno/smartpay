package types

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/guidogimeno/smartpay/pkg/scrapper"
	"github.com/shopspring/decimal"
)

const (
	months = 12
)

type Payment struct {
	Amount            float64
	Installments      int
	InstallmentAmount float64
	InterestRate      float32
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
		return errors.New("Amount must be greater than 0")
	}

	if p.InterestRate < 0 {
		return errors.New("Interest rate must be positive")
	}

	if p.InstallmentAmount < 0 {
		return errors.New("Installment amount must be positive")
	}

	if p.Installments < 1 {
		return errors.New("Number of installments must be greater than 0")
	}

	return nil
}

func (p *Payment) Analysis() ([]*Analysis, error) {
	var analysis []*Analysis
	ratables := []scrapper.Ratable{
		&scrapper.MercadoPago{},
		&scrapper.BCRA{},
	}

	var (
		resultsCh = make(chan *Analysis)
		errCh     = make(chan error)
		wg        = sync.WaitGroup{}
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, r := range ratables {
		wg.Add(1)
		go p.doTheMath(ctx, r, resultsCh, errCh, &wg)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
		close(errCh)
	}()

	for {
		select {
		case a, ok := <-resultsCh:
			if !ok {
				return analysis, nil
			}
			log.Print("a", a)
			analysis = append(analysis, a)
		case err := <-errCh:
			cancel()
			return nil, err
		}
	}
}

func (p *Payment) doTheMath(
	ctx context.Context,
	ratable scrapper.Ratable,
	resultsCh chan<- *Analysis,
	errCh chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	rate, err := ratable.Rate(ctx, time.Now())
	if err != nil {
		errCh <- errors.New("Failed to fetch BCRA data")
		return
	}

	tna := decimal.NewFromFloat(rate.Index).Div(decimal.NewFromInt(100))
	tnm := tna.Div(decimal.NewFromInt(months))
	installmentWithInterest := p.installmentWithInterest()

	savings := p.DecimalAmount().Copy()
	for i := 1; i <= p.Installments; i++ {
		fixedDepositInterest := savings.Mul(tnm)
		savings = savings.Add(fixedDepositInterest).Sub(installmentWithInterest)
	}

	resultsCh <- &Analysis{
		Entity:  rate.Source,
		Savings: savings.RoundDown(0).InexactFloat64(),
		Index:   rate.Index,
	}
}

func (p *Payment) installmentWithInterest() decimal.Decimal {
	if p.InstallmentAmount > 0 {
		return decimal.NewFromFloat(p.InstallmentAmount)
	}

	installments := decimal.NewFromInt(int64(p.Installments))
	interestRate := decimal.NewFromFloat32(p.InterestRate)

	installmentAmount := p.DecimalAmount().Div(installments)
	interest := installmentAmount.Mul(interestRate)

	return installmentAmount.Add(interest)
}
