package accrual

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

type storage interface {
	GetAllNotProcessedOrders(ctx context.Context) ([]entity.Order, error)

	SetOrderStatusAndAccrual(ctx context.Context, order entity.Order) error
}
type config interface {
	GetAccrualAddress() string
}

type Accrual struct {
	storage
	config
}

type accrualRespose struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

func New(config config, storage storage) (*Accrual, error) {
	return &Accrual{
		storage,
		config,
	}, nil
}

func (a *Accrual) Run(ctx context.Context) {
	accrualAddress := a.config.GetAccrualAddress()
	restClient := resty.New()
	restClient.SetDebug(true)
	updateTicker := time.NewTicker(1 * time.Second)
	defer updateTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-updateTicker.C:
			orders, _ := a.GetAllNotProcessedOrders(ctx)
			for _, order := range orders {
				accrualResp := accrualRespose{}
				resp, err := restClient.R().
					EnableTrace().
					SetResult(&accrualResp).
					Get(fmt.Sprintf("%s/api/orders/%d", accrualAddress, order.Number))
				if err != nil {
					fmt.Println(err)
				}
				if resp.StatusCode() == http.StatusOK {
					newOrder := entity.Order{
						Number:  order.Number,
						Status:  accrualResp.Status,
						Accrual: accrualResp.Accrual,
					}
					_ = a.SetOrderStatusAndAccrual(ctx, newOrder)
				}
			}
		}
	}
}
