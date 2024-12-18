package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/transport/accrual/model"
)

const (
	goodsEndpoint  = "/api/goods"
	ordersEndpoint = "/api/orders"
)

type HTTPClient struct {
	http *http.Client
	cfg  config.Config
	log  logger.Logger
}

func NewHTTPClient(cfg config.Config, log logger.Logger) *HTTPClient {
	return &HTTPClient{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		cfg: cfg,
		log: log,
	}
}

func (c *HTTPClient) RegisterNewGoods(ctx context.Context) error {
	fn := "accrual.RegisterNewGoods"

	goods := model.Goods{
		Match:      randomBrand(),
		Reward:     randomReward(),
		RewardType: model.Percent,
	}

	goodsData, err := json.Marshal(goods)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't marshal goods data", fn), "error", err.Error())
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL()+goodsEndpoint,
		bytes.NewReader(goodsData),
	)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't create request", fn), "error", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't do request", fn), "error", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code - %d, actual - %d", http.StatusOK, resp.StatusCode)
	}

	return nil
}

func (c *HTTPClient) RegisterNewOrder(ctx context.Context, orderID string) error {
	fn := "accrual.RegisterNewOrder"

	details := model.OrderDetails{
		Order: orderID,
		Goods: []model.OrderGoods{
			{
				Description: randomProductType() + " " + randomBrand(),
				Price:       randomPrice(),
			},
		},
	}

	detailsData, err := json.Marshal(details)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't marshal order details data", fn), "error", err.Error())
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL()+ordersEndpoint,
		bytes.NewReader(detailsData),
	)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't create request", fn), "error", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't do request", fn), "error", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("expected status code - %d, actual - %d", http.StatusAccepted, resp.StatusCode)
	}

	return nil
}

func (c *HTTPClient) UpdateOrderStatus(ctx context.Context, orderID string) (model.OrderStatus, error) {
	fn := "accrual.UpdateOrderStatus"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL()+ordersEndpoint+"/"+orderID,
		nil,
	)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't create request", fn), "error", err.Error())
		return model.OrderStatus{}, err
	}

	req.Header.Set("Content-Type", "plain/text")

	resp, err := c.http.Do(req)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't do request", fn), "error", err.Error())
		return model.OrderStatus{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		c.log.Error(fmt.Sprintf("%s: response status %s", fn, resp.Status), "err", "too many requests")
		after := resp.Header.Get("Retry-After")
		duration, err := time.ParseDuration(after)
		if err != nil {
			c.log.Error(fmt.Sprintf("%s: can't parse retry duration", fn), "error", err.Error())
			return model.OrderStatus{}, err
		}
		return model.OrderStatus{}, &ManyRequestError{RetryAfter: duration}
	}

	var orderStatus model.OrderStatus
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error(fmt.Sprintf("%s: can't read response body", fn), "error", err.Error())
		return model.OrderStatus{}, err
	}

	if err = json.Unmarshal(b, &orderStatus); err != nil {
		c.log.Error(fmt.Sprintf("%s: can't unmarshal response body", fn), "error", err.Error())
		return model.OrderStatus{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return model.OrderStatus{}, fmt.Errorf("expected status code - %d, actual - %d", http.StatusOK, resp.StatusCode)
	}

	return orderStatus, nil
}

func (c *HTTPClient) baseURL() string {
	return "http://localhost" + c.cfg.AccrualAddress()
}

func randomBrand() string {
	var brands = []string{"LG", "Apple", "Samsung"}
	return brands[rand.Intn(len(brands))]
}

func randomReward() float64 {
	var rewards = []float64{5, 10, 20}
	return rewards[rand.Intn(len(rewards))]
}

func randomProductType() string {
	var productTypes = []string{"TV", "Phone", "Monitor", "Camera"}
	return productTypes[rand.Intn(len(productTypes))]
}

func randomPrice() float64 {
	var price = []float64{300, 600, 1200, 3000}
	return price[rand.Intn(len(price))]
}
