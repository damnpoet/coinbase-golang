package price

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/text/currency"

	"github.com/damnpoet/coinbase-golang/config"
	"github.com/damnpoet/coinbase-golang/models"
	"github.com/damnpoet/coinbase-golang/net"
)

type resource struct {
	Price models.Price `json:"data"`
}

type Repository interface {
	Get(src string, dest currency.Unit) (*models.Price, error)
}

type repository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewRepository(config config.Reader, gateway net.Gateway) Repository {
	return &repository{
		config:  config,
		gateway: gateway,
	}
}

func (repo repository) Get(src string, dest currency.Unit) (*models.Price, error) {
	response, err := repo.gateway.Get(fmt.Sprintf("%s/prices/%s-%s/buy", repo.config.APIEndpoint(), src, dest.String()))
	if err != nil {
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(response.Body)

	var r resource
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON response from server: %s", err.Error())
	}

	return &r.Price, err
}
