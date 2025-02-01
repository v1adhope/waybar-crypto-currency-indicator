# Usage

The module uses [CoinMarketCap API](https://coinmarketcap.com/api/)

Each run consumes 1 credit, for a free account 10,000 per month is given.

Three environmental variables are required for full work:
- `CRYPTO_CURRENCY_API` - api token
- `CRYPTO_CURRENCY_FV` - coin for bar (USD by default)
- `CRYPTO_CURRENCY_WL` - coins for drop down list (BTC,ETH,DOT by default)

```bash
# To search for coin IDs, consumes 1 credit `*required curl and jq`
curl -H "X-CMC_PRO_API_KEY: <your_api>" -H "Accept: application/json" -d "symbol=<ticker_with_a_comma_separating>" -G https://pro-api.coinmarketcap.com/v1/cryptocurrency/map | jq .
```

## Sample

Updated every 30 minutes.

```
"custom/crypto-currency": {
    "max-length": 15,
    "return-type": "json",
    "format": "{}",
    "exec": "$HOME/.config/waybar/scripts/crypto-currency",
    "interval": 1800
}
```

# Preview

![preview](/assets/preview.png)
