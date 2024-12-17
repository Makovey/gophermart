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
	goodsEndpoint = "/api/goods"
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
	
	port := c.cfg.AccrualAddress()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://localhost%s%s", port, goodsEndpoint),
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

	c.log.Info(fmt.Sprintf("%s: response status %s", fn, resp.Status))

	return nil
}

func (c *HTTPClient) RegisterNewOrder(ctx context.Context, orderID string) error {
	return nil
}

func (c *HTTPClient) UpdateOrderStatus(ctx context.Context, orderID string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:8085/api/orders/12345", nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("response Status:", resp.Status)

	return nil
}

func randomBrand() string {
	var brands = []string{"LG", "Apple", "Samsung"}
	return brands[rand.Intn(len(brands))]
}

func randomReward() float64 {
	var rewards = []float64{5, 10, 20}
	return rewards[rand.Intn(len(rewards))]
}
