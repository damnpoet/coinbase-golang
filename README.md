# Coinbase Wallet Golang Client

This is a golang client library for the [Coinbase Wallet API v2](https://developers.coinbase.com/api/v2).

_Important:_ As this library is targeted for newer API v2, it requires v2 permissions (i.e. `wallet:accounts:read`).

### Installation

```go
$ go get github.com/damnpoet/coinbase-golang
```

## Authentication

### API Key (HMAC Client)

```ruby
require 'coinbase/wallet'
client = Coinbase::Wallet::Client.new(api_key: <api key>, api_secret: <api secret>)
```

```go
import "github.com/damnpoet/coinbase-golang"

c := coinbase.NewClient("<api key>", "<api secret>")
```

