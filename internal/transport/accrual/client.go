package accrual

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	
	"github.com/Makovey/gophermart/internal/config"
)

type HTTPClient struct {
	http   *http.Client
	config config.Config
}

func NewHTTPClient(cfg config.Config) *HTTPClient {
	return &HTTPClient{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		config: cfg,
	}
}

func (c *HTTPClient) RegisterNewGoods() error {
	return nil
}

func (c *HTTPClient) RegisterNewOrder(orderID string) error {
	return nil
}

func (c *HTTPClient) UpdateOrderStatus(orderID string) error {
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
