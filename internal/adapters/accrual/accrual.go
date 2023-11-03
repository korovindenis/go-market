package accrual

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

const maxWorker = 10

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
		storage: storage,
		config:  config,
	}, nil
}

func (a *Accrual) Run(ctx context.Context) {
	accrualAddress := a.config.GetAccrualAddress()
	restClient := resty.New()

	notProcessedOrdersCH := make(chan entity.Order, maxWorker)
	defer close(notProcessedOrdersCH)

	updateTicker := time.NewTicker(1 * time.Second)
	defer updateTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-updateTicker.C:
			orders, _ := a.GetAllNotProcessedOrders(ctx)
			for _, order := range orders {
				notProcessedOrdersCH <- order

				go a.worker(ctx, restClient, accrualAddress, notProcessedOrdersCH)
			}
		}
	}
}

func (a *Accrual) worker(ctx context.Context, restClient *resty.Client, accrualAddress string, orderCh <-chan entity.Order) {
	for {
		select {
		case <-ctx.Done():
			return
		case order := <-orderCh:
			accrualResp := accrualRespose{}
			resp, _ := restClient.R().
				EnableTrace().
				SetResult(&accrualResp).
				Get(fmt.Sprintf("%s/api/orders/%s", accrualAddress, order.Number))

			if resp.StatusCode() == http.StatusTooManyRequests {
				retryAfterHeader := resp.Header().Get("Retry-After")
				if retryAfterHeader != "" {
					secondsToWait, _ := strconv.Atoi(retryAfterHeader)
					time.Sleep(time.Duration(secondsToWait) * time.Second)
				}

				continue
			}

			if resp.StatusCode() == http.StatusOK {
				newOrder := entity.Order{
					Number:  order.Number,
					Status:  accrualResp.Status,
					Accrual: accrualResp.Accrual,
				}
				_ = a.SetOrderStatusAndAccrual(ctx, newOrder)
			}

			return
		}
	}
}
