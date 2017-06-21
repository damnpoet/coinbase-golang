package net

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/damnpoet/coinbase-golang/config"
)

const (
	DefaultDialTimeout = 5 * time.Second
)

type apiErrorHandler func(statusCode int, body []byte) error

type Gateway struct {
	errHandler   apiErrorHandler
	trustedCerts []tls.Certificate
	config       config.Reader
	transport    *http.Transport
	DialTimeout  time.Duration
}

func NewGateway(config config.Reader) *Gateway {
	return &Gateway{
		config: config,
	}
}

func (gateway Gateway) Get(url string) (*http.Response, error) {
	request, err := gateway.NewRequest("GET", url, gateway.config.Secret(), gateway.config.Key(), nil)
	if err != nil {
		return nil, err
	}

	return gateway.doRequestAndHandlerError(request)
}

func (gateway Gateway) NewRequest(method, path, secret, key string, body io.ReadSeeker) (*http.Request, error) {
	request, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %s", err.Error())
	}

	return gateway.newRequest(request, secret, key), nil
}

func (gateway Gateway) newRequest(request *http.Request, secret, key string) *http.Request {
	if secret != "" && key != "" {
		t := time.Now()
		timestamp := strconv.Itoa(int(t.Unix()))

		var bodyBuffer []byte
		if request.Body != nil {
			bodyBuffer, _ = ioutil.ReadAll(request.Body)
			request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
		}

		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte(timestamp))
		mac.Write([]byte(request.Method))
		mac.Write([]byte(request.URL.Path))
		mac.Write(bodyBuffer)

		request.Header.Set("CB-ACCESS-KEY", key)
		request.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
		request.Header.Set("CB-ACCESS-SIGN", base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set("Connection", "close")
	request.Header.Set("content-type", "application/json")

	return request
}

func (gateway Gateway) doRequestAndHandlerError(request *http.Request) (*http.Response, error) {
	rawResponse, err := gateway.doRequest(request)
	if err != nil {
		return rawResponse, WrapNetworkErrors(request.URL.Host, err)
	}

	if rawResponse.StatusCode > 299 {
		defer rawResponse.Body.Close()
		jsonBytes, _ := ioutil.ReadAll(rawResponse.Body)
		rawResponse.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBytes))
		err = gateway.errHandler(rawResponse.StatusCode, jsonBytes)
	}

	return rawResponse, err
}

func (gateway Gateway) doRequest(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error

	if gateway.transport == nil {
		makeHTTPTransport(&gateway)
	}

	httpClient := NewHTTPClient(gateway.transport, NewRequestDumper())

	httpClient.DumpRequest(request)

	for i := 0; i < 3; i++ {
		response, err = httpClient.Do(request)
		if response == nil && err != nil {
			continue
		} else {
			break
		}
	}

	if err != nil {
		return response, err
	}

	httpClient.DumpResponse(response)
	return response, err
}

func makeHTTPTransport(gateway *Gateway) {
	gateway.transport = &http.Transport{
		Dial: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			Timeout:   gateway.DialTimeout,
		}).Dial,
		TLSClientConfig: NewTLSConfig(gateway.trustedCerts, gateway.config.IsSSLDisabled()),
		Proxy:           http.ProxyFromEnvironment,
	}
}

func dialTimeout(envDialTimeout string) time.Duration {
	dialTimeout := DefaultDialTimeout
	if timeout, err := strconv.Atoi(envDialTimeout); err == nil {
		dialTimeout = time.Duration(timeout) * time.Second
	}
	return dialTimeout
}

func (gateway *Gateway) SetTrustedCerts(certificates []tls.Certificate) {
	gateway.trustedCerts = certificates
	makeHTTPTransport(gateway)
}
