package config

type Reader interface {
	APIEndpoint() string
	IsSSLDisabled() bool
	Secret() string
	Key() string
}

type config struct {
	apiEndpoint string
	secret      string
	key         string
	sslDisabled bool
}

func New(apiEndpoint, secret, key string) Reader {
	return &config{
		apiEndpoint: apiEndpoint,
		secret:      secret,
		key:         key,
	}
}

func (c *config) APIEndpoint() string {
	return c.apiEndpoint
}

func (c *config) IsSSLDisabled() bool {
	return c.sslDisabled
}

func (c *config) Secret() string {
	return c.secret
}

func (c *config) Key() string {
	return c.key
}
