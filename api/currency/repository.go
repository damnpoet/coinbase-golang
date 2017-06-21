package currency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/damnpoet/coinbase-golang/config"
	"github.com/damnpoet/coinbase-golang/models"
	"github.com/damnpoet/coinbase-golang/net"
)

type resource struct {
	Currencies []models.Currency `json:"data"`
}

type Repository interface {
	List() ([]models.Currency, error)
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

func (repo repository) List() ([]models.Currency, error) {
	response, err := repo.gateway.Get(fmt.Sprintf("%s/currencies", repo.config.APIEndpoint()))
	if err != nil {
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(response.Body)

	var r resource
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON response from server: %s", err.Error())
	}

	return r.Currencies, err
}
