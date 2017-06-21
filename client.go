package coinbase

import (
	currencyAPI "github.com/damnpoet/coinbase-golang/api/currency"
	"github.com/damnpoet/coinbase-golang/api/price"
	"github.com/damnpoet/coinbase-golang/config"
	"github.com/damnpoet/coinbase-golang/models"
	"github.com/damnpoet/coinbase-golang/net"

	"golang.org/x/text/currency"
)

const (
	APIv2Url = "https://api.coinbase.com/v2"

	Bitcoin  = "BTC"
	Ethereum = "ETH"
	Litecoin = "LTC"
)

type Client interface {
	GetCurrencies() ([]models.Currency, error)
	GetBuyPrice(src string, dest currency.Unit) (*models.Price, error)
}

type hmacClient struct {
	gateway *net.Gateway
	conf    config.Reader
}

// NewClient creates a new HMAC Client
func NewClient(apiEndpoint, key, secret string) Client {
	c := config.New(apiEndpoint, key, secret)

	return &hmacClient{
		conf:    c,
		gateway: net.NewGateway(c),
	}
}

func (c *hmacClient) GetCurrencies() ([]models.Currency, error) {
	repo := currencyAPI.NewRepository(c.conf, *c.gateway)
	return repo.List()
}

func (c *hmacClient) GetBuyPrice(src string, dest currency.Unit) (*models.Price, error) {
	repo := price.NewRepository(c.conf, *c.gateway)
	return repo.Get(src, dest)
}
